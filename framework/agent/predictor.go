package agent

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"

	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/downloadmanager"
	"github.com/rai-project/utils"
	context "golang.org/x/net/context"
)

type Predictor struct {
	Host string
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

// isBase64 tests a string to determine if it is a base64 or not.
func isBase64(toTest string) bool {
	_, err := base64.StdEncoding.DecodeString(toTest)
	return err == nil
}

func (p *Predictor) InputReaderCloser(ctx context.Context, req *dl.PredictRequest) (io.ReadCloser, error) {
	var data string

	if bts, err := base64.StdEncoding.DecodeString(string(req.Data)); err == nil {
		data = string(bts)
	} else {
		data = string(req.Data)
	}

	pp.Println("url = ", data)

	if utils.IsURL(data) {
		targetDir, err := UploadDir()
		if err != nil {
			return nil, err
		}
		path, err := downloadmanager.Download(data, targetDir)
		if err != nil {
			return nil, err
		}
		f, err := os.Open(path)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open file %v", path)
		}
		return f, nil
	}
	if data != "" {
		return ioutil.NopCloser(bytes.NewBufferString(data)), nil
	}
	return nil, errors.Errorf("invalid input data to InputReaderCloser")
}

func (p *Predictor) PublishInRegistery() error {
	return p.Base.PublishInPredictor(p.Host, "predictor")
}
