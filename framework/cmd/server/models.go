package server

import (
	"github.com/rai-project/config"
	"github.com/rai-project/dlframework"
)

var (
	models []dlframework.ModelManifest
)

func init() {
	config.AfterInit(func() {
		models = framework.Models()
	})
}
