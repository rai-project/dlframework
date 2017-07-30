package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModelManifests(t *testing.T) {
	manifests, err := models.allmanifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetTensorflowModelManifests(t *testing.T) {
	manifests, err := models.manifests("tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetInceptionModelManifests(t *testing.T) {
	manifests, err := models.manifests("tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
	ms, err := models.filter(manifests, "inception", "*")
	assert.NoError(t, err)
	assert.NotEmpty(t, ms)
}
