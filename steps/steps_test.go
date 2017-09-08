package steps

import (
	"testing"

	"github.com/rai-project/image/types"
	"github.com/rai-project/pipeline"
	"github.com/stretchr/testify/assert"
)

func TestURLReadImage(t *testing.T) {
	imgURLs := []string{
		"https://jpeg.org/images/jpeg-home.jpg",
		"https://www.google.com/images/branding/googlelogo/1x/googlelogo_color_272x92dp.png",
	}

	input := make(chan string)
	for _, url := range imgURLs {
		input <- url
	}

	output := pipeline.New().
		Then(NewReadURL()).
		Then(NewReadImage()).
		Run(input)

	for out := range output {
		assert.NotEmpty(t, out)
		assert.IsType(t, types.RGBImage{}, out)
	}
}
