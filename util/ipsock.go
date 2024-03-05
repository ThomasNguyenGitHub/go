package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

const (
	HeaderRealIP       = "X-Real-IP"
	HeaderForwardedFor = "X-Forwarded-For"
)

// GetIP gets the IP address from an incoming HTTP request.
func GetIP(r *http.Request) (string, error) {
	// Get IP from the X-Real-IP header
	realIP := r.Header.Get(HeaderRealIP)
	if len(net.ParseIP(realIP)) > 0 {
		return realIP, nil
	}

	// Get IP from X-Forwarded-For header
	forwardedFor := r.Header.Get(HeaderForwardedFor)
	ips := strings.Split(forwardedFor, ",")
	for _, ip := range ips {
		if len(net.ParseIP(ip)) > 0 {
			return ip, nil
		}
	}

	// Get IP from RemoteAddr
	remoteAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	if len(net.ParseIP(remoteAddr)) > 0 {
		return remoteAddr, nil
	}
	return "", fmt.Errorf("no valid IP address found")
}
