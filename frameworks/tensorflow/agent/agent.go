package agent

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rai-project/config"
	rgrpc "github.com/rai-project/grpc"
	"github.com/rai-project/uuid"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"

	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/downloadmanager"
	common "github.com/rai-project/dlframework/frameworks/common/agent"
	tf "github.com/rai-project/dlframework/frameworks/tensorflow"
	"github.com/rai-project/dlframework/frameworks/tensorflow/predict"
)

type registryServer struct {
	common.Registry
}

type predictorServer struct {
	common.Predictor
}

func (p *predictorServer) Predict(ctx context.Context, req *dl.PredictRequest) (*dl.PredictResponse, error) {
	framework, err := dl.FindFramework(req.GetFrameworkName() + ":" + req.GetFrameworkVersion())
	if err != nil {
		return nil, err
	}
	model, err := framework.FindModel(req.GetModelName() + ":" + req.GetModelVersion())
	if err != nil {
		return nil, err
	}

	predictor, err := predict.New(model)
	if err != nil {
		return nil, err
	}
	defer predictor.Close()
	if err := predictor.Download(); err != nil {
		return nil, err
	}

	var reader io.Reader
	if req.GetUrl() != "" {
		targetDir := filepath.Join(config.App.TempDir, "uploads")
		path, err := downloadmanager.Download(req.GetUrl(), targetDir)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file %v", path)
		}
		defer f.Close()
		reader = f
	} else if req.GetData() != nil {
		reader = bytes.NewBuffer(req.GetData())
	} else {
		return nil, errors.New("invalid input")
	}
	probs, err := predictor.Predict(reader)
	if err != nil {
		return nil, err
	}

	probs.Sort()

	if req.GetLimit() != 0 {
		trunc := probs.Take(int(req.GetLimit()))
		probs = &trunc
	}

	return &dl.PredictResponse{
		Id:       uuid.NewV4(),
		Features: *probs,
		Error:    nil,
	}, nil
}

func RegisterRegistryServer() *grpc.Server {
	var grpcServer *grpc.Server
	grpcServer = rgrpc.NewServer(dl.RegistryServiceDescription)
	svr := &registryServer{
		Registry: common.Registry{
			Base: common.Base{
				Framework: tf.FrameworkManifest,
			},
		},
	}
	svr.PublishInRegistery()
	dl.RegisterRegistryServer(grpcServer, svr)
	return grpcServer
}

func RegisterPredictorServer() *grpc.Server {
	var grpcServer *grpc.Server
	grpcServer = rgrpc.NewServer(dl.PredictorServiceDescription)
	svr := &predictorServer{
		Predictor: common.Predictor{
			Base: common.Base{
				Framework: tf.FrameworkManifest,
			},
		},
	}
	svr.PublishInRegistery()
	dl.RegisterPredictorServer(grpcServer, svr)
	return grpcServer
}
