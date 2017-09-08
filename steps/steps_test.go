package steps

import (
	"bytes"
	"testing"

	"golang.org/x/net/context"

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

	ctx := context.Background()
	output := pipeline.New(ctx).
		Then(NewReadURL()).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		assert.IsType(t, &bytes.Buffer{}, out)
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

	ctx := context.Background()
	output := pipeline.New(ctx).
		Then(NewReadURL()).
		Then(NewReadImage()).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		assert.IsType(t, &types.RGBImage{}, out)
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

	ctx := context.Background()
	output := pipeline.New(ctx).
		Then(NewReadURL()).
		Then(NewReadImage()).
		Then(NewPreprocessImage(predict.PreprocessOptions{})).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		assert.IsType(t, []float32{}, out)
	}
}
