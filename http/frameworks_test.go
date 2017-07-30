package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFrameworkManifests(t *testing.T) {
	manifests, err := frameworks.manifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetFrameworkManifestsFilter(t *testing.T) {
	manifests, err := frameworks.manifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
	fs, err := frameworks.filter(manifests, "tensorflow", "latest")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "tensorflow", "*")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "Tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "tensorflow", "1.2.x")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "tensorflow", "1.x")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = frameworks.filter(manifests, "tensorflow", ">=1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
}
