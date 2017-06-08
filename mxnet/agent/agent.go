package agent

import (
	"encoding/json"
	"errors"

	"github.com/levigross/grequests"
	"github.com/rai-project/dlframework/mxnet"
	"github.com/rai-project/grpc"
	context "golang.org/x/net/context"
)

type server struct{}

func (s *server) InferURL(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferURLServer) error {
	panic(errors.New("*server.InferURL not implemented"))
}

func (s *server) InferBytes(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferBytesServer) error {
	panic(errors.New("*server.InferBytes not implemented"))
}

func (s *server) GetModelGraph(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Graph, error) {
	model, err := mxnet.GetModelInformation(m.GetName())
	if err != nil {
		return nil, err
	}
	graphURL := model.GetGraphUrl()
	if graphURL == "" {
		return nil, errors.New("empty graph url")
	}
	resp, err := grequests.Get(graphURL)
	if err != nil {
		return nil, err
	}
	g := new(mxnet.Model_Graph)
	err := json.Unmarshal(resp.String(), g)
	return g, err
}

func (s *server) GetModelInformation(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Information, error) {
	model, err := mxnet.GetModelInformation(m.GetName())
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

func Register() {
	grpcServer := grpc.NewServer()
	mxnet.RegisterMXNetServer(grpcServer, &server{})
}

func NewServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	mxnet.RegisterMXNetServer(grpcServer, &server{})
	return grpcServer
}
