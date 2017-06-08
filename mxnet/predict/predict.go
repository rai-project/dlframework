package predict

import (
	"bufio"
	"errors"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
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
	Predict(data interface{}) ([]float32, error)
}

type ImagePredictor struct {
	model    mxnet.Model_Information
	modelDir string
}

type Feature struct {
	idx  int
	name string
	prob float32
}

func NewImagePredictor(model mxnet.Model_Information, targetDir string) (Predictor, error) {
	return &ImagePredictor{
		model:    model,
		modelDir: targetDir,
	}, nil
}

func (p *ImagePredictor) GetGraphPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+"-graph.json")
}

func (p *ImagePredictor) GetSymbolPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+"-symbol.json")
}

func (p *ImagePredictor) GetWeightsPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+"-weights.params")
}

func (p *ImagePredictor) GetFeaturesPath() string {
	return filepath.Join(p.modelDir, p.model.GetName()+".features")
}

func (p *ImagePredictor) Preprocess(input interface{}) (interface{}, error) {
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

func (p *ImagePredictor) Predict(input string) ([]Feature, error) {
	img, err := imgio.Open(input)
	if err != nil {
		return nil, err
	}

	_, ok := img.([]float32)
	if !ok {
		return nil, errors.New("expecting an flattened float32 array input")
	}

	symbol := ioutil.ReadFile(model.GetSymbolPath())
	params := ioutil.ReadFile(model.GetWeightsPath())

	modelInput := model.GetInput()
	modelInputShape := modelInput.GetDimensions()

	var features []string
	f, _ := os.Open(model.GetFeaturesPath)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		features = append(features, line)
	}

	p, err := gomxnet.CreatePredictor(symbol,
		params,
		gomxnet.Device{mxnet.CPU_DEVICE, 0},
		[]mxnet.InputNode{{Key: "data", Shape: modelInputShape}},
	)
	if err != nil {
		return nil, err
	}
	defer p.Free()

	resized := transform.Resize(img, int(modelInputShape[2]), int(modelInputShape[3]), transform.Linear)
	res, err := Preprocess(resized)
	if err != nil {
		return nil, err
	}

	if err := p.SetInput("data", res); err != nil {
		return nil, err
	}

	if err := p.Forward(); err != nil {
		return nil, err
	}

	probs, err := p.GetOutput(0)
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
		ret[i].name = features[out.Idxs[i]]
	}

	return ret, nil
}
