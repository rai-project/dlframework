package mxnet

import (
	"encoding/json"
	"testing"

	rice "github.com/GeertJohan/go.rice"
	"github.com/stretchr/testify/assert"
)

var (
	box                 = rice.MustFindBox("_fixtures")
	inceptionSymbolJSON = box.MustBytes("Inception-BN-symbol.json")
	caffenetSymbolJSON  = box.MustBytes("caffenet-symbol.json")
)

func TestUnmarshalGraph(t *testing.T) {
	var g Model_Graph
	err := json.Unmarshal(caffenetSymbolJSON, &g)
	assert.NoError(t, err)
	assert.NotEmpty(t, g)

	s, err := json.MarshalIndent(g, "", "  ")
	assert.NoError(t, err)
	assert.NotEmpty(t, s)
	// t.Log(string(s))

	dg, err := g.ToDotGraph()
	assert.NoError(t, err)
	assert.NotNil(t, dg)
	assert.NotEmpty(t, dg)

	t.Log(dg.String())

}
