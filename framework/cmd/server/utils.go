package server

import (
	"net"
)

func getTracerHostAddress() string {
	host, _, err := net.SplitHostPort(tracerAddress)
	if err != nil {
		return tracerAddress
	}
	return host
}
