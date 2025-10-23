package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIPAddr(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if netIp := net.ParseIP(ip); netIp.IsPrivate() || netIp.IsLoopback() {
		if header_ip := r.Header.Get("X-Forwarded-For"); net.ParseIP(header_ip) != nil {
			ip = header_ip
		}

	}
	return ip
}
