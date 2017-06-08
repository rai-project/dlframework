package predict

import (
	"bufio"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/pkg/errors"
	"github.com/rai-project/dlframework/mxnet"
	gomxnet "github.com/songtianyi/go-mxnet-predictor/mxnet"
	"github.com/songtianyi/go-mxnet-predictor/utils"
)

type Predictor interface {
	// Downloads the features / symbol file / weights
	Download() error
	// Preprocess the data
	Preprocess(data interface{}) (interface{}, error)
	// Returns the features
	Predict(data interface{}) ([]Feature, error)

	io.Closer
}

type ImagePredictor struct {
	model     mxnet.Model_Information
	modelDir  string
	features  []string
	predictor *gomxnet.Predictor
}

type Feature struct {
	idx  int
	name string
	prob float32
}

func NewImagePredictor(model mxnet.Model_Information, targetDir string) (*ImagePredictor, error) {
	return &ImagePredictor{
		model:    model,
		modelDir: targetDir,
	}, nil
}

func (p *ImagePredictor) GetGraphPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+"-graph.json")
}

func (p *ImagePredictor) GetWeightsPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+"-weights.params")
}

func (p *ImagePredictor) GetFeaturesPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+".features")
}

func (p *ImagePredictor) Preprocess(input interface{}) ([]float32, error) {
	img, ok := input.(image.Image)
	if !ok {
		return nil, errors.New("expecting an image input")
	}

	model := p.model
	meanImage := model.GetMeanImage()
	if len(meanImage) == 0 {
		meanImage = []float32{0, 0, 0}
	}

	b := img.Bounds()
	h := b.Max.Y - b.Min.Y // image height
	w := b.Max.X - b.Min.X // image width

	res := make([]float32, 3*h*w)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, _ := img.At(x+b.Min.X, y+b.Min.Y).RGBA()
			res[y*w+x] = float32(r>>8) - meanImage[0]
			res[w*h+y*w+x] = float32(g>>8) - meanImage[1]
			res[2*w*h+y*w+x] = float32(b>>8) - meanImage[2]
		}
	}
	return res, nil
}

func (p *ImagePredictor) getPredictor() error {
	model := p.model

	symbol, err := ioutil.ReadFile(p.GetGraphPath())
	if err != nil {
		return errors.Wrapf(err, "cannot read %s", p.GetGraphPath())
	}
	params, err := ioutil.ReadFile(p.GetWeightsPath())
	if err != nil {
		return errors.Wrapf(err, "cannot read %s", p.GetWeightsPath())
	}

	var features []string
	f, err := os.Open(p.GetFeaturesPath())
	if err != nil {
		return errors.Wrapf(err, "cannot read %s", p.GetFeaturesPath())
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		features = append(features, line)
	}

	p.features = features

	modelInput := model.GetInput()
	t := modelInput.GetDimensions()

	modelInputShape := make([]uint32, len(t))
	for i := range t {
		modelInputShape[i] = uint32(t[i])
	}

	pred, err := gomxnet.CreatePredictor(symbol,
		params,
		gomxnet.Device{gomxnet.CPU_DEVICE, 0},
		[]gomxnet.InputNode{{Key: "data", Shape: modelInputShape}},
	)
	if err != nil {
		return err
	}
	p.predictor = pred

	return nil
}

func (p *ImagePredictor) Predict(input interface{}) ([]Feature, error) {
	path, ok := input.(string)
	if !ok {
		return nil, errors.New("expecting a path to predict")
	}
	img, err := imgio.Open(path)
	if err != nil {
		return nil, err
	}

	if p.predictor == nil {
		if err := p.getPredictor(); err != nil {
			return nil, err
		}
	}

	modelInput := p.model.GetInput()
	t := modelInput.GetDimensions()

	resized := transform.Resize(img, int(t[2]), int(t[3]), transform.Linear)
	res, err := p.Preprocess(resized)
	if err != nil {
		return nil, err
	}

	if err := p.predictor.SetInput("data", res); err != nil {
		return nil, err
	}

	if err := p.predictor.Forward(); err != nil {
		return nil, err
	}

	probs, err := p.predictor.GetOutput(0)
	if err != nil {
		return nil, err
	}

	idxs := make([]int, len(probs))
	for i := range probs {
		idxs[i] = i
	}
	out := utils.ArgSort{Args: probs, Idxs: idxs}
	sort.Sort(out)

	ret := make([]Feature, len(probs))
	for i := range probs {
		ret[i].prob = out.Args[i]
		ret[i].idx = out.Idxs[i]
		ret[i].name = p.features[out.Idxs[i]]
	}

	return ret, nil
}

func (p *ImagePredictor) Close() error {
	if p.predictor != nil {
		p.predictor.Free()
	}
	return nil
}
