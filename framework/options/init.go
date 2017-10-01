package options

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
)

var (
	log = logger.New().WithField("pkg", "dlframework/predictor/options")
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "dlframework/predictor/options")
	})
}
