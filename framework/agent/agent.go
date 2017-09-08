package agent

import (
	"fmt"
	"sync"

	// _ "github.com/rai-project/dldataset/vision"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/utils"
	"github.com/rai-project/uuid"
	context "golang.org/x/net/context"
	"golang.org/x/sync/syncmap"

	"google.golang.org/grpc"

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
}

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
		predictor: predictor,
		options:   options,
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

	predictor, err = p.predictor.Load(ctx, *model)
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4()
	p.loadedPredictors.Store(id, predictor)

	return &dl.Predictor{Id: id}, nil
}

func (p *Agent) getLoadedPredictor(ctx context.Context, id string) (predictor.ImagePredictor, error) {

	val, ok := p.loadedPredictors.Load(id)
	if !ok {
		return nil, errors.Errorf("predictor %v was not found", id)
	}

	predictor, ok := val.(predict.ImagePredictor)
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

	predictor.Close()

	loadedPredictors.Delete(id)

	return &dl.PredictorCloseResponse{}, nil
}

// Image method receives a stream of urls and runs
// the predictor on all the urls. The
//
// The result is a prediction feature stream for each url.
func (p *Agent) URLs(context.Context, *URLsRequest) (*dl.FeaturesResponse, error) {
	ctx := svr.Context()

	predictorId := req.GetPredictor().GetId()

	predictor, err := p.getLoadedPredictor(ctx, predictorId)
	if err != nil {
		return nil, err
	}

	_, model, err := predictor.Info()
	if err != nil {
		return nil, err
	}

	preprocessOptions, err := predictor.PreprocessOptions()
	if err != nil {
		return nil, err
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range req.GetUrls() {
			input <- *url
		}
	}()

	output := pipeline.New(ctx).
		Then(steps.NewReadURL()).
		Then(steps.NewReadImage()).
		Then(steps.NewPreprocessImage(preprocessOptions)).
		Then(steps.NewImagePredict(predictor)).
		Run(input)

	var wg sync.WaitGroup
  var mutex sync.Mutex

  res := &dl.FeaturesResponse{}

	for out := range output {
		if err, ok := out.(error); ok {
			return nil, err
		}
		o, ok := out.(steps.IDer)
		if !ok {
			return errors.Errorf("expecting an ider type, but got %v", o)
		}

		features, ok := o.GetData().([]*Feature)
		if !ok {
			return errors.Errorf("expecting a []*Feature type, but got %v", o.GetData())
		}

		wg.Add(1)
		go func() {
      defer wg.Done()
      mutex.Lock()
      defer mutex.Unlock()
			res.Responses = append(res.Responses, &dl.FeatureResponse{
				Id:        uuid.NewV4(),
				InputId:   o.GetId(),
				RequestId: "todo-request-id",
				Features:  features,
			}))
		}()
	}
	wg.Wait()

	return res, nil
}

// Image method receives a list base64 encoded images and runs
// the predictor on all the images.
//
// The result is a prediction feature stream for each image.
func (p *Agent) Images(req *dl.ImagesRequest) (*dl.FeaturesResponse, error) {
	ctx := svr.Context()

    predictorId := req.GetPredictor().GetId()

    predictor, err := p.getLoadedPredictor(ctx, predictorId)
    if err != nil {
      return nil, err
    }

    _, model, err := predictor.Info()
    if err != nil {
      return nil, err
    }

    preprocessOptions, err := predictor.PreprocessOptions()
    if err != nil {
      return nil, err
    }

    input := make(chan interface{})
    go func() {
      defer close(input)
      for _, img := range req.GetImages() {
        input <- *img
      }
    }()

    output := pipeline.New(ctx).
      Then(steps.NewReadImage()).
      Then(steps.NewPreprocessImage(preprocessOptions)).
      Then(steps.NewImagePredict(predictor)).
      Run(input)

    var wg sync.WaitGroup
    var mutex sync.Mutex

    res := &dl.FeaturesResponse{}

    for out := range output {
      if err, ok := out.(error); ok {
        return nil, err
      }
      o, ok := out.(steps.IDer)
      if !ok {
        return errors.Errorf("expecting an ider type, but got %v", o)
      }

      features, ok := o.GetData().([]*Feature)
      if !ok {
        return errors.Errorf("expecting a []*Feature type, but got %v", o.GetData())
      }

      wg.Add(1)
      go func() {
        defer wg.Done()
        mutex.Lock()
        defer mutex.Unlock()
        res.Responses = append(res.Responses, &dl.FeatureResponse{
          Id:        uuid.NewV4(),
          InputId:   o.GetId(),
          RequestId: "todo-request-id",
          Features:  features,
        }))
      }()
    }
    wg.Wait()

    return res, nil
}

// Dataset method receives a single dataset and runs
// the predictor on all elements of the dataset.
//
// The result is a prediction feature stream.
func (p *Agent) Dataset(req *dl.DatasetRequest) (*dl.FeaturesResponse, error) {
	return nil, nil
}

// Clear method clears the internal cache of the predictors
func (p *Agent) Reset(ctx context.Context, req *dl.ResetRequest) (*dl.ResetResponse, error) {
	return nil, nil
}

// func (p *Agent) Predict(ctx context.Context, req *dl.PredictRequest) (*dl.PredictResponse, error) {
// 	_, model, err := p.FindFrameworkModel(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	predictor, err := p.predictor.Load(ctx, *model)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer predictor.Close()
// 	if err := predictor.Download(ctx); err != nil {
// 		return nil, err
// 	}

// 	reader, err := p.ReadInput(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer reader.Close()

// 	var img image.Image

// 	func() {
// 		if span, newCtx := opentracing.StartSpanFromContext(ctx, "DecodeImage"); span != nil {
// 			ctx = newCtx
// 			defer span.Finish()
// 		}
// 		img, _, err = image.Decode(reader)
// 	}()
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "unable to read input as image")
// 	}

// 	data, err := predictor.Preprocess(ctx, img)
// 	if err != nil {
// 		return nil, err
// 	}

// 	probs, err := predictor.Predict(ctx, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	probs.Sort()

// 	if req.GetLimit() != 0 {
// 		trunc := probs.Take(int(req.GetLimit()))
// 		probs = &trunc
// 	}

// 	return &dl.PredictResponse{
// 		Id:       uuid.NewV4(),
// 		Features: *probs,
// 		Error:    nil,
// 	}, nil
// }

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

// func (p *Agent) ReadInput(ctx context.Context, req *dl.PredictRequest) (io.ReadCloser, error) {
// 	if span, newCtx := opentracing.StartSpanFromContext(ctx, "ReadInput"); span != nil {
// 		ctx = newCtx
// 		defer span.Finish()
// 	}

// 	data := tryBase64Decode(req.Data)

// 	if data == "" {
// 		return nil, errors.Errorf("invalid empty input data to ReadInput")
// 	}

// 	if strings.HasPrefix(data, "dataset://") {
// 		pth := strings.TrimPrefix(data, "dataset://")
// 		sep := strings.SplitAfterN(pth, "/", 3)
// 		if len(sep) != 3 {
// 			return nil, errors.Errorf("the dataset path %s is not formatted correctly expected datasets://category/name/file_path", data)
// 		}
// 		category := sep[0]
// 		name := sep[1]
// 		rest := sep[2]

// 		dataset, err := dldataset.Get(category, name)
// 		if err != nil {
// 			return nil, err
// 		}

// 		err = dataset.Download(ctx)
// 		if err != nil {
// 			return nil, err
// 		}

// 		label, err := dataset.Get(ctx, rest)
// 		if err != nil {
// 			return nil, err
// 		}
// 		iface, err := label.Data()
// 		if err != nil {
// 			return nil, err
// 		}
// 		if reader, ok := iface.(io.Reader); ok {
// 			return ioutil.NopCloser(reader), nil
// 		}
// 		pp.Println("TODO.. we need to still support images as output...")
// 		return nil, errors.New("unhandeled dataset input...")
// 	}

// 	if utils.IsURL(data) {
// 		targetDir, err := UploadDir()
// 		if err != nil {
// 			return nil, err
// 		}
// 		path, err := downloadmanager.DownloadInto(data, targetDir)
// 		if err != nil {
// 			return nil, err
// 		}
// 		f, err := os.Open(path)
// 		if err != nil {
// 			return nil, errors.Wrapf(err, "failed to open file %v", path)
// 		}
// 		return f, nil
// 	}

// 	return ioutil.NopCloser(bytes.NewBufferString(data)), nil
// }

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
