package generator

import (
	"fmt"
	"testing"

	"github.com/rai-project/dlframework/mxnet/models"
	"github.com/stretchr/testify/assert"
)

func TestGenerator(t *testing.T) {
	model, err := models.Get("locationnet")
	assert.NoError(t, err)
	assert.NotEmpty(t, model)

	generator, err := New(model)
	assert.NoError(t, err)
	assert.NotEmpty(t, generator)

	f := generator.Generate()

	fmt.Printf("%#v", f)
}
