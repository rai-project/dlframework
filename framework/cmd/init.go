package cmd

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/sirupsen/logrus"
)

var (
	log       *logrus.Entry = logrus.New().WithField("pkg", "dlframework/framework/cmd")
	debugging               = false
)

func init() {
	log.Level = logrus.DebugLevel
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "dlframework/framework/cmd")
	})
}
