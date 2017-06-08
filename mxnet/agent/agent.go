package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"

	"google.golang.org/grpc"

	"github.com/levigross/grequests"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework/mxnet"
	"github.com/rai-project/dlframework/mxnet/predict"
	rgrpc "github.com/rai-project/grpc"
	context "golang.org/x/net/context"
)

type server struct{}

func (s *server) InferURL(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferURLServer) error {
	resp, err := grequests.Get(m.GetUrl(), nil)
	if err != nil {
		return err
	}
	defer resp.Close()
	m.Data = resp.Bytes()
	return s.InferBytes(m, m1)
}

func (s *server) InferBytes(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferBytesServer) error {
	model, err := mxnet.GetModelInformation(m.GetModelName())
	if err != nil {
		return err
	}
	predictor, err := predict.NewImagePredictor(model, config.App.TempDir)
	if err != nil {
		m1.Send(&mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		})
		return err
	}
	defer predictor.Close()

	if err := predictor.Download(); err != nil {
		m1.Send(&mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		})
		return err
	}

	img, _, err := image.Decode(bytes.NewBuffer(m.GetData()))
	if err != nil {
		m1.Send(&mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		})
		return err
	}

	pre, err := predictor.Preprocess(img)
	if err != nil {
		m1.Send(&mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		})
		return err
	}
	features, err := predictor.Predict(pre)
	if err != nil {
		m1.Send(&mxnet.MXNetInferenceResponse{
			Error: &mxnet.ErrorStatus{
				Ok:      false,
				Message: err.Error(),
			},
		})
		return err
	}
	m1.Send(&mxnet.MXNetInferenceResponse{
		Id:       m.GetId(),
		Features: features,
		Error: &mxnet.ErrorStatus{
			Ok: true,
		},
	})
	return nil
}

func (s *server) GetModelGraph(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Graph, error) {
	model, err := mxnet.GetModelInformation(m.GetModelName())
	if err != nil {
		return nil, err
	}
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
