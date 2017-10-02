package agent

import (
	"os"
	"strconv"

	"github.com/facebookgo/freeport"
	"github.com/pkg/errors"
	tr "github.com/rai-project/tracer"
	"github.com/rai-project/utils"
)

type Options struct {
	host   string
	port   int
	tracer tr.Tracer
}

type Option func(o *Options) *Options

func WithTracer(tracer tr.Tracer) Option {
	return func(o *Options) *Options {
		o.tracer = tracer
		return o
	}
}

func WithHost(host string) Option {
	return func(o *Options) *Options {
		o.host = host
		return o
	}
}

func WithPort(port int) Option {
	return func(o *Options) *Options {
		o.port = port
		return o
	}
}

func WithPortString(port string) Option {
	return func(o *Options) *Options {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.WithError(err).WithField("port", port).Error("unable to parse port")
			return o
		}
		o.port = p
		return o
	}
}

func getPort() (int, error) {
	port, found := os.LookupEnv("PORT")
	if !found {
		return freeport.Get()
	}
	return strconv.Atoi(port)
}

func getHost() (string, error) {
	return utils.GetLocalIp()
}

func NewOptions(opts ...Option) (*Options, error) {
	host, err := getHost()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent host address")
	}
	port, err := getPort()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent listen port")
	}
	options := &Options{
		host: host,
		port: port,
	}
	for _, o := range opts {
		o(options)
	}
	return options, nil
}
