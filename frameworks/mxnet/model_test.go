package mxnet

import (
	"os"
	"testing"

	rice "github.com/GeertJohan/go.rice"
	"github.com/k0kubun/pp"
	"github.com/rai-project/dlframework"
	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v2"
)

var (
	builtinModelsBox = rice.MustFindBox("builtin_models")
)

func TestUnmarshalModel(t *testing.T) {
	builtinModelsBox.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var model dlframework.ModelManifest
		data := builtinModelsBox.MustBytes(path)
		err = yaml.Unmarshal(data, &model)
		assert.NoError(t, err)
		assert.NotEmpty(t, model)
		pp.Println(model)
		return err
	})

}
