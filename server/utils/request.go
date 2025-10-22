package utils

import (
	"net"
	"net/http"
	"strings"
)

func GetIPAddr(r *http.Request) string {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	if net.ParseIP(ip).IsPrivate() {
		ip = r.Header.Get("X-Forwarded-For")
	}
	return ip
}
