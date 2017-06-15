package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"google.golang.org/grpc"

	"github.com/levigross/grequests"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework/frameworks/mxnet"
	"github.com/rai-project/dlframework/frameworks/mxnet/predict"
	rgrpc "github.com/rai-project/grpc"
	context "golang.org/x/net/context"

	gocache "github.com/patrickmn/go-cache"
)

type server struct{}

func (s *server) InferURL(ctx context.Context, m *mxnet.MXNetInferenceRequest) (*mxnet.MXNetInferenceResponse, error) {
	resp, err := grequests.Get(m.GetUrl(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Close()
	m.Data = resp.Bytes()
	return s.InferBytes(ctx, m)
}

func (s *server) InferBytes(ctx context.Context, m *mxnet.MXNetInferenceRequest) (*mxnet.MXNetInferenceResponse, error) {

	model, err := mxnet.GetModelInformation(m.GetModelName())
	if err != nil {
		return nil, err
	}
	predictor, err := predict.NewImagePredictor(model, config.App.TempDir)
	if err != nil {
		return &mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		}, err
	}
	defer predictor.Close()

	if err := predictor.Download(); err != nil {
		return &mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		}, err
	}

	img, _, err := image.Decode(bytes.NewBuffer(m.GetData()))
	if err != nil {
		return &mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		}, err
	}

	pre, err := predictor.Preprocess(img)
	if err != nil {
		return &mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		}, err
	}
	features, err := predictor.Predict(pre)
	if err != nil {
		return &mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		}, err
	}
	log.Debug("infered features....")
	return &mxnet.MXNetInferenceResponse{
		Id:       m.GetId(),
		Features: features,
		Error: &mxnet.ErrorStatus{
			Ok: true,
		},
	}, nil
}

func (s *server) GetModelGraph(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Graph, error) {
	model, err := mxnet.GetModelInformation(m.GetModelName())
	if err != nil {
		return nil, err
	}

	cacheKey := "graph/" + model.GetName()
	if val, found := cache.Get(cacheKey); found {
		if g, ok := val.(*mxnet.Model_Graph); ok {
			return g, nil
		}
	}

	log.WithField("url", model.GetGraphUrl()).WithField("cacheKey", cacheKey).Debug("downloading model graph")

	graphURL := model.GetGraphUrl()
	if graphURL == "" {
		return nil, errors.New("empty graph url")
	}
	resp, err := grequests.Get(graphURL, nil)
	if err != nil {
		return nil, err
	}
	g := new(mxnet.Model_Graph)
	err = json.Unmarshal(resp.Bytes(), g)
	if err == nil {
		cache.Set(cacheKey, g, gocache.DefaultExpiration)
	}
	return g, err
}

func (s *server) GetModelInformation(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Information, error) {
	model, err := mxnet.GetModelInformation(m.GetModelName())
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (s *server) GetModelInformations(ctx context.Context, n *mxnet.Null) (*mxnet.ModelInformations, error) {
	names := mxnet.ModelNames()
	models := make([]*mxnet.Model_Information, len(names))
	for ii, name := range names {
		model, err := mxnet.GetModelInformation(name)
		if err != nil {
			return nil, err
		}
		models[ii] = &model
	}
	return &mxnet.ModelInformations{
		Info: models,
	}, nil
}

func Register() *grpc.Server {
	var grpcServer *grpc.Server
	grpcServer = rgrpc.NewServer(mxnet.ServiceDescription)
	mxnet.RegisterMXNetServer(grpcServer, &server{})
	return grpcServer
}
