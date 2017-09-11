package registryquery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFrameworkManifests(t *testing.T) {
	manifests, err := Frameworks.manifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
}

func TestGetFrameworkManifestsFilterManifests(t *testing.T) {
	manifests, err := Frameworks.manifests()
	assert.NoError(t, err)
	assert.NotEmpty(t, manifests)
	fs, err := Frameworks.FilterManifests(manifests, "tensorflow", "latest")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "tensorflow", "*")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "Tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "tensorflow", "1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "tensorflow", "1.2.x")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "tensorflow", "1.x")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
	fs, err = Frameworks.FilterManifests(manifests, "tensorflow", ">=1.2")
	assert.NoError(t, err)
	assert.NotEmpty(t, fs)
}
