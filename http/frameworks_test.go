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
