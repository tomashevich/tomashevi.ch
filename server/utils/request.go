package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIPAddr(r *http.Request, checkXFF bool) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]

	if !checkXFF {
		return ip
	}

	if netIp := net.ParseIP(ip); netIp.IsLoopback() || netIp.IsPrivate() {
		xff := r.Header.Get("X-Forwarded-For")
		if xff == "" {
			return ip
		}

		xffIp := strings.Split(xff, ",")[0]
		if net.ParseIP(xffIp) == nil {
			return ip
		}

		ip = xffIp
	}

	return ip
}
