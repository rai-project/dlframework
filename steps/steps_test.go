package steps

import (
	"bytes"
	"testing"

	"github.com/rai-project/dlframework/framework/predict"
	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	"github.com/stretchr/testify/assert"
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
			input <- url
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
			input <- url
		}
	}()

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage(predict.PreprocessOptions{
			MeanImage: []float32{0, 0, 0},
			Size:      []int{224, 224},
			Scale:     1.0,
			ColorMode: types.RGBMode,
		})).
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
			input <- url
		}
	}()

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage(predict.PreprocessOptions{
			MeanImage: []float32{0, 0, 0},
			Size:      []int{224, 224},
			Scale:     1.0,
			ColorMode: types.RGBMode,
		})).
		Then(NewPreprocessImage(predict.PreprocessOptions{})).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		v, ok := out.(IDer)
		assert.True(t, ok)
		assert.IsType(t, []float32{}, v.GetData())
	}
}
