package server

import (
	"net"
	"strings"
)

func getTracerHostAddress(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return strings.TrimPrefix(addr, "http://")
	}
	return strings.TrimPrefix(host, "http://")
}
