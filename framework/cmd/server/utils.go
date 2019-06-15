package server

import (
	"net"
	"strings"
)

func getTracerHostAddress(addr string) string {
	trimPrefix := func(s string) string {
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, "http://") {
			return strings.TrimPrefix(s, "http://")
		}
		if strings.HasPrefix(s, "https://") {
			return strings.TrimPrefix(s, "https://")
		}
		return s
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return trimPrefix(addr)
	}
	return trimPrefix(host)
}
