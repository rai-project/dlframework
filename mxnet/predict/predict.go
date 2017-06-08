package predict

import (
	"errors"
	"image"
	"path/filepath"

	"github.com/rai-project/dlframework/mxnet"
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

func NewImagePredictor(model mxnet.Model_Information, targetDir string) (Predictor, error) {
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

func (p *ImagePredictor) Predict(input interface{}) ([]float32, error) {
	data, ok := input.([]float32)
	if !ok {
		return nil, errors.New("expecting an flattened float32 array input")
	}

	_ = data
	// model := p.model
	// modelInput := model.GetInput()
	// modelInputShape := modelInput.GetDimensions()
	// p, err := mxnet.CreatePredictor(symbol,
	// 	params,
	// 	mxnet.Device{mxnet.CPU_DEVICE, 0},
	// 	[]mxnet.InputNode{{Key: "data", Shape: inputShape}},
	// )
	// if err != nil {
	// 	return 0, "", err
	// }
	// defer p.Free()
	// // etc...
	return nil, nil
}
