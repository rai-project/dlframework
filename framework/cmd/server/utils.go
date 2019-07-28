package server

import (
	"net"
	"strings"

	externalip "github.com/glendc/go-external-ip"
)

func getTracerServerAddress(addr string) string {
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
	host = trimPrefix(host)
	if host == "localhost" {
		consensus := externalip.DefaultConsensus(nil, nil)
		ip, err := consensus.ExternalIP()
		if err != nil {
			return ""
		}
		return ip.String()
	}
	return host
}
