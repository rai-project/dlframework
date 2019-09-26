package server

import (
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	_ "github.com/rai-project/tracer/jaeger"
	_ "github.com/rai-project/tracer/noop"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "dlframework/framework/cmd/server")
	})
}
