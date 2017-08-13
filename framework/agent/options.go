package agent

import (
	"os"
	"strconv"

	"github.com/facebookgo/freeport"
	"github.com/pkg/errors"
	"github.com/rai-project/utils"
)

type Options struct {
	host string
	port int
}

type Option func(o *Options) *Options

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

func NewOptions() (*Options, error) {
	host, err := getHost()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent host address")
	}
	port, err := getPort()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get agent listen port")
	}
	return &Options{
		host: host,
		port: port,
	}, nil
}
