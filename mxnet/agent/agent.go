package agent

import (
	"errors"

	"github.com/rai-project/dlframework/mxnet"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) GetModel(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model, error) {
	panic(errors.New("*server.GetModel not implemented"))
}

func (s *server) GetModelGraph(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Graph, error) {
	panic(errors.New("*server.GetModelGraph not implemented"))
}

func (s *server) GetModelInformation(ctx context.Context, m *mxnet.MXNetModelInformationRequest) (*mxnet.Model_Information, error) {
	model, err := mxnet.GetModelInformation(m.GetName())
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (s *server) Infer(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferServer) error {
	panic(errors.New("*server.Infer not implemented"))
}

func (s *server) InferBytes(m *mxnet.MXNetInferenceRequest, m1 mxnet.MXNet_InferBytesServer) error {
	panic(errors.New("*server.InferBytes not implemented"))
}

func (s *server) GetModels(ctx context.Context, n *mxnet.Null) (*mxnet.ModelInformations, error) {
	names := mxnet.ModelNames()
	models := make([]*Model_Information, len(names))
	for ii, name := range names {
		model, err := mxnet.GetModelInformation(name)
		if err != nil {
			return nil, err
		}
		models[ii] = model
	}
	return &mxnet.ModelInformations{
		Info: models,
	}, nil
}

func Register() {
	grpcServer := grpc.NewServer()
	mxnet.RegisterMXNetServer(grpcServer, &server{})
}
