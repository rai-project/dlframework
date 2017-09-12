package agent

import (
	"fmt"
	"reflect"
	"sync"

	// _ "github.com/rai-project/dldataset/vision"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/utils"
	"github.com/rai-project/uuid"
	context "golang.org/x/net/context"
	"golang.org/x/sync/syncmap"

	"google.golang.org/grpc"

	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	rgrpc "github.com/rai-project/grpc"
	"github.com/rai-project/pipeline"
	"github.com/rai-project/registry"

	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/dlframework/steps"
)

type Agent struct {
	base
	loadedPredictors syncmap.Map
	predictor        predict.Predictor
	options          *Options
	channelBuffer    int
}

var (
	DefaultChannelBuffer = 1000
)

func New(predictor predict.Predictor, opts ...Option) (*Agent, error) {
	options, err := NewOptions()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(options)
	}
	framework, _, err := predictor.Info()
	if err != nil {
		return nil, err
	}
	return &Agent{
		base: base{
			Framework: framework,
		},
		predictor:     predictor,
		options:       options,
		channelBuffer: DefaultChannelBuffer,
	}, nil
}

// Opens a predictor and returns an id where the predictor
// is accessible. The id can be used to perform inference
// requests.
func (p *Agent) Open(ctx context.Context, req *dl.PredictorOpenRequest) (*dl.Predictor, error) {
	_, model, err := p.FindFrameworkModel(ctx, req)
	if err != nil {
		return nil, err
	}

	predictor, err := p.predictor.Load(ctx, *model)
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4()
	p.loadedPredictors.Store(id, predictor)

	return &dl.Predictor{Id: id}, nil
}

func (p *Agent) getLoadedPredictor(ctx context.Context, id string) (predict.Predictor, error) {
	val, ok := p.loadedPredictors.Load(id)
	if !ok {
		return nil, errors.Errorf("predictor %v was not found", id)
	}

	predictor, ok := val.(predict.Predictor)
	if !ok {
		return nil, errors.Errorf("predictor %v is not a valid image predictor", predictor)
	}

	return predictor, nil
}

// Close a predictor clear it's memory.
func (p *Agent) Close(ctx context.Context, req *dl.Predictor) (*dl.PredictorCloseResponse, error) {
	id := req.Id
	predictor, err := p.getLoadedPredictor(ctx, id)
	if err != nil {
		return nil, err
	}

	predictor.Reset(ctx)

	p.loadedPredictors.Delete(id)

	return &dl.PredictorCloseResponse{}, nil
}

func (p *Agent) toFeaturesResponse(output <-chan interface{}, options *dl.PredictionOptions) (*dl.FeaturesResponse, error) {

	var wg sync.WaitGroup
	var mutex sync.Mutex

	res := &dl.FeaturesResponse{}

	if options == nil {
		options = &dl.PredictionOptions{
			RequestId: "request-id-not-found",
		}
	}

	for out := range output {
		if err, ok := out.(error); ok {
			return nil, err
		}

		var features []*dl.Feature

		inputId := "undefined-input-id"

		switch o := out.(type) {
		case steps.IDer:
			inputId = o.GetId()
			data := o.GetData()
			switch e := data.(type) {
			case []*dl.Feature:
				features = e
			case dl.Features:
				features = []*dl.Feature(e)
			default:
				return nil, errors.Errorf("expecting a []*Feature type, but got %v", data)
			}
		case []*dl.Feature:
			features = o
		case dl.Features:
			features = []*dl.Feature(o)
		default:
			return nil, errors.Errorf("expecting an ider or []*Feature type, but got %v", reflect.TypeOf(out))
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			mutex.Lock()
			defer mutex.Unlock()
			res.Responses = append(res.Responses, &dl.FeatureResponse{
				Id:        uuid.NewV4(),
				InputId:   inputId,
				RequestId: options.GetRequestId(),
				Features:  features,
			})
		}()
	}
	wg.Wait()

	return res, nil
}

// Image method receives a stream of urls and runs
// the predictor on all the urls. The
//
// The result is a prediction feature list for each url.
func (p *Agent) URLs(ctx context.Context, req *dl.URLsRequest) (*dl.FeaturesResponse, error) {

	if req.GetPredictor() == nil {
		return nil, errors.New("request does not have a valid predictor set")
	}

	predictorId := req.GetPredictor().GetId()
	if predictorId == "" {
		return nil, errors.New("predictor id cannot be an empty string")
	}

	predictor, err := p.getLoadedPredictor(ctx, predictorId)
	if err != nil {
		return nil, err
	}

	preprocessOptions, err := predictor.PreprocessOptions(ctx)
	if err != nil {
		return nil, err
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range req.GetUrls() {
			input <- url
		}
	}()

	output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(p.channelBuffer)).
		Then(steps.NewReadURL()).
		Then(steps.NewReadImage(preprocessOptions)).
		Then(steps.NewPreprocessImage(preprocessOptions)).
		Then(steps.NewPredictImage(predictor)).
		Run(input)

	return p.toFeaturesResponse(output, req.GetOptions())
}

// Image method receives a stream of urls and runs
// the predictor on all the urls. The
//
// The result is a prediction feature stream for each url.
func (p *Agent) URLsStream(req *dl.URLsRequest, svr dl.Predict_URLsStreamServer) error {
	return nil
}

// Image method receives a list base64 encoded images and runs
// the predictor on all the images.
//
// The result is a prediction feature list for each image.
func (p *Agent) Images(ctx context.Context, req *dl.ImagesRequest) (*dl.FeaturesResponse, error) {

	if req.GetPredictor() == nil {
		return nil, errors.New("request does not have a valid predictor set")
	}
	predictorId := req.GetPredictor().GetId()
	if predictorId == "" {
		return nil, errors.New("predictor id cannot be an empty string")
	}

	predictor, err := p.getLoadedPredictor(ctx, predictorId)
	if err != nil {
		return nil, err
	}

	preprocessOptions, err := predictor.PreprocessOptions(ctx)
	if err != nil {
		return nil, err
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, img := range req.GetImages() {
			input <- img
		}
	}()

	output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(p.channelBuffer)).
		Then(steps.NewReadImage(preprocessOptions)).
		Then(steps.NewPreprocessImage(preprocessOptions)).
		Then(steps.NewPredictImage(predictor)).
		Run(input)

	return p.toFeaturesResponse(output, req.GetOptions())
}

// Image method receives a list base64 encoded images and runs
// the predictor on all the images.
//
// The result is a prediction feature stream for each image.
func (p *Agent) ImagesStream(req *dl.ImagesRequest, svr dl.Predict_ImagesStreamServer) error {
	return nil
}

// Dataset method receives a single dataset and runs
// the predictor on all elements of the dataset.
//
// The result is a prediction feature list.
func (p *Agent) Dataset(ctx context.Context, req *dl.DatasetRequest) (*dl.FeaturesResponse, error) {

	if req.GetPredictor() == nil {
		return nil, errors.New("request does not have a valid predictor set")
	}
	predictorId := req.GetPredictor().GetId()
	if predictorId == "" {
		return nil, errors.New("predictor id cannot be an empty string")
	}

	predictor, err := p.getLoadedPredictor(ctx, predictorId)
	if err != nil {
		return nil, err
	}

	preprocessOptions, err := predictor.PreprocessOptions(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetDataset() == nil {
		return nil, errors.New("invalid empty dataset parameter in request")
	}

	dataset, err := dldataset.Get(req.Dataset.GetCategory(), req.Dataset.GetName())
	if err != nil {
		return nil, err
	}

	dataset, err = dataset.New(ctx)
	if err != nil {
		return nil, err
	}

	defer dataset.Close()

	if err := dataset.Download(ctx); err != nil {
		return nil, err
	}

	elems, err := dataset.List(ctx)
	if err != nil {
		return nil, err
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, e := range elems {
			input <- e
		}
	}()

	output := pipeline.New(pipeline.Context(ctx), pipeline.ChannelBuffer(p.channelBuffer)).
		Then(steps.NewGetDataset(dataset)).
		Then(steps.NewReadImage(preprocessOptions)).
		Then(steps.NewPreprocessImage(preprocessOptions)).
		Then(steps.NewPredictImage(predictor)).
		Run(input)

	return p.toFeaturesResponse(output, req.GetOptions())
}

// Dataset method receives a single dataset and runs
// the predictor on all elements of the dataset.
//
// The result is a prediction feature stream.
func (p *Agent) DatasetStream(req *dl.DatasetRequest, svr dl.Predict_DatasetStreamServer) error {
	return nil
}

// Clear method clears the internal cache of the predictors
func (p *Agent) Reset(ctx context.Context, req *dl.ResetRequest) (*dl.ResetResponse, error) {
	return nil, nil
}

func (p *Agent) FindFrameworkModel(ctx context.Context, req *dl.PredictorOpenRequest) (*dl.FrameworkManifest, *dl.ModelManifest, error) {
	framework, err := dl.FindFramework(req.GetFrameworkName() + ":" + req.GetFrameworkVersion())
	if err != nil {
		return nil, nil, err
	}
	model, err := framework.FindModel(req.GetModelName() + ":" + req.GetModelVersion())
	if err != nil {
		return nil, nil, err
	}

	return framework, model, nil
}

func (p *Agent) RegisterManifests() (*grpc.Server, error) {
	log.Info("populating registry")

	var grpcServer *grpc.Server
	grpcServer = rgrpc.NewServer(dl.RegistryServiceDescription)
	svr := &Registry{
		base: base{
			Framework: p.base.Framework,
		},
	}
	go func() {
		utils.Every(
			registry.Config.Timeout/2,
			func() {
				svr.PublishInRegistery()
			},
		)
	}()
	dl.RegisterRegistryServer(grpcServer, svr)
	return grpcServer, nil
}

func (p *Agent) RegisterPredictor() (*grpc.Server, error) {

	grpcServer := rgrpc.NewServer(dl.PredictServiceDescription)

	host := fmt.Sprintf("%s:%d", p.options.host, p.options.port)
	log.Info("registering predictor service at ", host)

	go func() {
		utils.Every(
			registry.Config.Timeout/2,
			func() {
				p.base.PublishInPredictor(host, "predictor")
			},
		)
	}()
	dl.RegisterPredictServer(grpcServer, p)
	return grpcServer, nil
}
