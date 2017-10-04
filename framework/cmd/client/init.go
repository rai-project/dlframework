package client

import (
	"os"
	"syscall"
	"time"

	"github.com/k0kubun/pp"
	shutdown "github.com/klauspost/shutdown2"
	"github.com/rai-project/config"
	"github.com/rai-project/logger"
	"github.com/rai-project/tracer"
	_ "github.com/rai-project/tracer/jaeger"
	_ "github.com/rai-project/tracer/noop"
	_ "github.com/rai-project/tracer/zipkin"
	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Entry
)

func init() {
	config.AfterInit(func() {
		log = logger.New().WithField("pkg", "dlframework/framework/cmd/client")
	})

	shutdown.OnSignal(0, os.Interrupt, syscall.SIGTERM)
	shutdown.SetTimeout(time.Second * 1)
	shutdown.SecondFn(func() {
		pp.Println("🛑 shutting down!!")
		tracer.Close()
	})
}
