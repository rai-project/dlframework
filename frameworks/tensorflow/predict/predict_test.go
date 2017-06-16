package predict

import (
	"os"
	"testing"

	"github.com/rai-project/config"
	tf "github.com/rai-project/dlframework/frameworks/tensorflow"
	"github.com/stretchr/testify/assert"
)

func XXXTestPredictLoad(t *testing.T) {
	framework := tf.FrameworkManifest
	model, err := framework.FindModel("vgg19:1.0")
	assert.NoError(t, err)
	assert.NotEmpty(t, model)

	predictor, err := New(model)
	assert.NoError(t, err)
	assert.NotEmpty(t, predictor)

	defer predictor.Close()

	imgPredictor, ok := predictor.(*ImagePredictor)
	assert.True(t, ok)

	assert.NotEmpty(t, imgPredictor.imageDimensions)
	assert.NotEmpty(t, imgPredictor.meanImage)

}

func TestPredictInference(t *testing.T) {
	framework := tf.FrameworkManifest
	model, err := framework.FindModel("inception:3.0")
	assert.NoError(t, err)
	assert.NotEmpty(t, model)

	predictor, err := New(model)
	assert.NoError(t, err)
	assert.NotEmpty(t, predictor)
	defer predictor.Close()

	err = predictor.Download()
	assert.NoError(t, err)
}

func TestMain(m *testing.M) {
	config.Init(
		config.AppName("carml"),
		config.DebugMode(true),
		config.VerboseMode(true),
	)
	os.Exit(m.Run())
}
