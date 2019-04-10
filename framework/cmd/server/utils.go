package server

import (
	"net"
)

func getTracerHostAddress(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	return host
}
