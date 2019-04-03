package registryquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModelManifests(t *testing.T) {
	manifests, err := Models.AllManifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetTensorFlowModelManifests(t *testing.T) {
	manifests, err := Models.Manifests("tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetInceptionModelManifests(t *testing.T) {
	manifests, err := Models.Manifests("tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
	ms, err := Models.FilterManifests(manifests, "inception", "*")
	assert.NoError(t, err)
	assert.NotEmpty(t, ms)
}
