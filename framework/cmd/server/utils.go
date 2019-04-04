package server

import (
	"net/url"

	"github.com/k0kubun/pp"
)

func getTracerHostAddress() string {
	u, err := url.Parse(tracerAddress)
	if err != nil {
		return tracerAddress
	}
	if u.Host == "" {
		return tracerAddress
	}
	pp.Println(u.Host)
	return u.Host
}
