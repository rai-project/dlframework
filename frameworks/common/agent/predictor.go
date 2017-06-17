package agent

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/downloadmanager"
	context "golang.org/x/net/context"
)

type Predictor struct {
	Base
}

func (p *Predictor) Predict(ctx context.Context, req *dl.PredictRequest) (*dl.PredictResponse, error) {
	return nil, nil
}

func (p *Predictor) FindFrameworkModel(ctx context.Context, req *dl.PredictRequest) (*dl.FrameworkManifest, *dl.ModelManifest, error) {
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

func (p *Predictor) InputReaderCloser(ctx context.Context, req *dl.PredictRequest) (io.ReadCloser, error) {
	if req.GetUrl() != "" {
		targetDir, err := UploadDir()
		if err != nil {
			return nil, err
		}
		path, err := downloadmanager.Download(req.GetUrl(), targetDir)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file %v", path)
		}
		return f, nil
	}
	if req.GetData() != nil {
		return ioutil.NopCloser(bytes.NewBuffer(req.GetData())), nil
	}
	return nil, errors.New("invalid input")
}

func (p *Predictor) PublishInRegistery() error {
	return p.Base.PublishInRegistery("predictor")
}
