package mxnet

import (
	"testing"

	"github.com/k0kubun/pp"
	"github.com/stretchr/testify/assert"
)

func TestBuiltinModel(t *testing.T) {
	info, err := GetModelInformation("nin")
	assert.NoError(t, err)
	assert.NotEmpty(t, info)
	pp.Println(info)
}
