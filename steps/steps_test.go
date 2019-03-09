package steps

import (
	"bytes"
	"os"
	"testing"

	"github.com/k0kubun/pp"

	"github.com/rai-project/config"
	dl "github.com/rai-project/dlframework"
	"github.com/rai-project/dlframework/framework/predictor"
	"github.com/rai-project/image"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	_ "github.com/rai-project/tracer/jaeger"
	"github.com/rai-project/uuid"
	"github.com/stretchr/testify/assert"
	"gorgonia.org/tensor"
)

func TestURLRead(t *testing.T) {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range imgURLs {
			input <- &dl.URLsRequest_URL{
				ID:   uuid.NewV4(),
				Data: url,
			}
		}
	}()

	output := pipeline.New().
		Then(NewReadURL()).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		v, ok := out.(IDer)
		assert.True(t, ok)
		assert.IsType(t, &bytes.Buffer{}, v.GetData())
	}
}

func TestURLReadImage(t *testing.T) {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range imgURLs {
			input <- &dl.URLsRequest_URL{
				ID:   uuid.NewV4(),
				Data: url,
			}
		}
	}()

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage(predictor.PreprocessOptions{})).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		v, ok := out.(IDer)
		assert.True(t, ok)
		assert.IsType(t, &types.RGBImage{}, v.GetData())
	}
}

func TestURLReadPreprocessImage(t *testing.T) {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range imgURLs {
			input <- &dl.URLsRequest_URL{
				ID:   uuid.NewV4(),
				Data: url,
			}
		}
	}()

	opts := predictor.PreprocessOptions{
		MeanImage: []float32{128, 100, 104},
		Size:      []int{224, 224},
		Scale:     1.0,
		ColorMode: types.RGBMode,
		Layout:    image.HWCLayout,
	}

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage(opts)).
		Then(NewPreprocessImage(opts)).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		v, ok := out.(IDer)
		assert.True(t, ok)
		assert.IsType(t, []float32{}, v.GetData())
		data := v.GetData().([]float32)
		pp.Println(data[10000])
	}
}

func TestURLReadPreprocessImageTensor(t *testing.T) {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}

	input := make(chan interface{})
	go func() {
		defer close(input)
		for _, url := range imgURLs {
			input <- &dl.URLsRequest_URL{
				ID:   uuid.NewV4(),
				Data: url,
			}
		}
	}()

	opts := predictor.PreprocessOptions{
		MeanImage: []float32{128, 100, 104},
		Size:      []int{224, 224},
		Scale:     1.0,
		ColorMode: types.RGBMode,
		Layout:    image.HWCLayout,
	}

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage(opts)).
		Then(NewPreprocessImageTensor(opts)).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		v, ok := out.(IDer)
		assert.True(t, ok)
		assert.IsType(t, &tensor.Dense{}, v.GetData())
		data := v.GetData().(*tensor.Dense).Data().([]float32)
		pp.Println(data[10000])
	}
}

func TestMain(m *testing.M) {
	config.Init(
		config.AppName("carml"),
		config.VerboseMode(true),
		config.DebugMode(true),
	)
	if false {
		pp.Println("keep")
	}
	os.Exit(m.Run())
}
