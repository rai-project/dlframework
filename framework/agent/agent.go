package agent

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/rai-project/dldataset"
	_ "github.com/rai-project/dldataset/vision"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/utils"
	context "golang.org/x/net/context"

	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"google.golang.org/grpc"

	rgrpc "github.com/rai-project/grpc"
	"github.com/rai-project/registry"

	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/uuid"
)

type Agent struct {
	Base
	predictor predict.Predictor
	options   *Options
}

func New(predictor predict.Predictor, opts ...Option) (*Agent, error) {
	options, err := NewOptions()
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		opt(options)
	}
	framework, err := predictor.GetFramework()
	if err != nil {
		return nil, err
	}
	return &Agent{
		Base: Base{
			Framework: framework,
		},
		predictor: predictor,
		options:   options,
	}, nil
}

func (p *Agent) Predict(ctx context.Context, req *dl.PredictRequest) (*dl.PredictResponse, error) {
	_, model, err := p.FindFrameworkModel(ctx, req)
	if err != nil {
		return nil, err
	}

	predictor, err := p.predictor.Load(ctx, *model)
	if err != nil {
		return nil, err
	}
	defer predictor.Close()
	if err := predictor.Download(ctx); err != nil {
		return nil, err
	}

	reader, err := p.ReadInput(ctx, req)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var img image.Image

	func() {
		if span, newCtx := opentracing.StartSpanFromContext(ctx, "DecodeImage"); span != nil {
			ctx = newCtx
			defer span.Finish()
		}
		img, _, err = image.Decode(reader)
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read input as image")
	}

	data, err := predictor.Preprocess(ctx, img)
	if err != nil {
		return nil, err
	}

	probs, err := predictor.Predict(ctx, data)
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

func (p *Agent) FindFrameworkModel(ctx context.Context, req *dl.PredictRequest) (*dl.FrameworkManifest, *dl.ModelManifest, error) {
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

func (p *Agent) ReadInput(ctx context.Context, req *dl.PredictRequest) (io.ReadCloser, error) {
	if span, newCtx := opentracing.StartSpanFromContext(ctx, "ReadInput"); span != nil {
		ctx = newCtx
		defer span.Finish()
	}

	data := tryBase64Decode(req.Data)

	if data == "" {
		return nil, errors.Errorf("invalid empty input data to ReadInput")
	}

	if strings.HasPrefix(data, "dataset://") {
		pth := strings.TrimPrefix(data, "dataset://")
		sep := strings.SplitAfterN(pth, "/", 3)
		category := sep[0]
		name := sep[1]
		rest := sep[2]

		dataset, err := dldataset.Get(category, name)
		if err != nil {
			return nil, err
		}

		err = dataset.Download(ctx)
		if err != nil {
			return nil, err
		}

		label, err := dataset.Get(ctx, rest)
		if err != nil {
			return nil, err
		}
		reader, err := label.Data()
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(reader), nil
	}

	if utils.IsURL(data) {
		targetDir, err := UploadDir()
		if err != nil {
			return nil, err
		}
		path, err := downloadmanager.DownloadInto(ctx, data, targetDir)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file %v", path)
		}
		return f, nil
	}

	return ioutil.NopCloser(bytes.NewBufferString(data)), nil
}

func (p *Agent) RegisterManifests() (*grpc.Server, error) {
	log.Info("populating registry")

	var grpcServer *grpc.Server
	grpcServer = rgrpc.NewServer(dl.RegistryServiceDescription)
	svr := &Registry{
		Base: Base{
			Framework: p.Base.Framework,
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

	grpcServer := rgrpc.NewServer(dl.PredictorServiceDescription)

	host := fmt.Sprintf("%s:%d", p.options.host, p.options.port)
	log.Info("registering predictor service at ", host)

	go func() {
		utils.Every(
			registry.Config.Timeout/2,
			func() {
				p.Base.PublishInPredictor(host, "predictor")
			},
		)
	}()
	dl.RegisterPredictorServer(grpcServer, p)
	return grpcServer, nil
}
