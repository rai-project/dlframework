package predict

import (
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"

	"github.com/rai-project/archive"
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/downloadmanager"
	common "github.com/rai-project/dlframework/frameworks/common/predict"
	"github.com/rai-project/utils"

	"github.com/Unknwon/com"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type ImagePredictor struct {
	*common.Base
	meanImage       []float64
	imageDimensions []int32
	tfGraph         *tf.Graph
	tfSession       *tf.Session
	workDir         string
}

func New(model dlframework.ModelManifest) (common.Predictor, error) {
	modelInputs := model.GetInputs()
	if len(modelInputs) != 1 {
		return nil, errors.New("number of inputs not supported")
	}
	firstInputType := modelInputs[0].GetType()
	if strings.ToLower(firstInputType) != "image" {
		return nil, errors.New("input type not supported")
	}
	return newImagePredictor(model)
}

func newImagePredictor(model dlframework.ModelManifest) (*ImagePredictor, error) {
	framework, err := model.ResolveFramework()
	if err != nil {
		return nil, err
	}

	cannonicalName, err := model.CanonicalName()
	if err != nil {
		return nil, err
	}
	workDir := filepath.Join(config.App.TempDir, strings.Replace(cannonicalName, ":", "_", -1))
	if err := os.MkdirAll(workDir, 0700); err != nil {
		return nil, errors.Wrapf(err, "failed to create model work directory %v", workDir)
	}

	ip := &ImagePredictor{
		Base: &common.Base{
			Framework: framework,
			Model:     model,
		},
		workDir: workDir,
	}

	if err := ip.setImageDimensions(); err != nil {
		return nil, err
	}

	if err := ip.setMeanImage(); err != nil {
		return nil, err
	}

	return ip, nil
}

func (p *ImagePredictor) makeSession() error {

	model := []byte("temporary")

	// Construct an in-memory graph from the serialized form.
	graph := tf.NewGraph()
	if err := graph.Import(model, ""); err != nil {
		return errors.Wrap(err, "unable to create tensorflow model graph")
	}

	// Create a session for inference over graph.
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return errors.Wrap(err, "unable to create tensorflow session")
	}

	p.tfGraph = graph
	p.tfSession = session

	return nil
}

func (p *ImagePredictor) Download() error {
	model := p.Model
	if model.Model.IsArchive {
		baseURL := model.Model.BaseUrl
		parsedURL, err := url.Parse(baseURL)
		if err != nil {
			return errors.Wrapf(err, "%v is not a valid url. unable to parse it", baseURL)
		}
		fileBaseName := filepath.Base(parsedURL.Path)
		_ = fileBaseName
		target := p.workDir

		err = downloadmanager.Download(baseURL, target)
		if err != nil {
			return errors.Wrapf(err, "failed to download model archive from %v", model.Model.BaseUrl)
		}
		targetFile, err := os.Open(target)
		if err != nil {
			return errors.Wrapf(err, "unable to open %v file", target)
		}
		defer targetFile.Close()
		unarchivedPath := filepath.Join(p.workDir, "model")
		archive.Unzip(targetFile, unarchivedPath)
	}
	return nil
}

func (p *ImagePredictor) Preprocess(data interface{}) (interface{}, error) {
	return nil, nil
}

func (p *ImagePredictor) Predict(data interface{}) ([]*dlframework.PredictionFeature, error) {

	if p.tfSession == nil {
		if err := p.makeSession(); err != nil {
			return nil, err
		}
	}

	session := p.tfSession
	graph := p.tfGraph

	var reader io.ReadCloser
	defer func() {
		if reader != nil {
			reader.Close()
		}
	}()

	if str, ok := data.(string); ok {
		if com.IsFile(str) {
			f, err := os.Open(str)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to open file from %v", str)
			}
			reader = f
		} else if utils.IsURL(str) {
			resp, err := grequests.Get(str, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to download data from %v", str)
			}
			reader = resp
		}
	}
	if rdr, ok := data.(io.Reader); ok {
		reader = ioutil.NopCloser(rdr)
	}

	tensor, err := p.makeTensorFromImage(reader)
	if err != nil {
		return nil, err
	}
	output, err := session.Run(
		map[tf.Output]*tf.Tensor{
			graph.Operation("input").Output(0): tensor,
		},
		[]tf.Output{
			graph.Operation("output").Output(0),
		},
		nil)
	if err != nil {
		log.Fatal(err)
	}
	// output[0].Value() is a vector containing probabilities of
	// labels for each image in the "batch". The batch size was 1.
	// Find the most probably label index.
	probabilities := output[0].Value().([][]float32)[0]

	res := make([]*dlframework.PredictionFeature, len(probabilities))
	for ii, prob := range probabilities {
		res[ii] = &dlframework.PredictionFeature{
			Index:       int64(ii),
			Probability: prob,
		}
	}

	return res, nil
}

func (p *ImagePredictor) Close() error {
	if p.tfSession != nil {
		p.tfSession.Close()
	}
	return nil
}
