package ipx

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIP(r *http.Request) string {
	ipSources := []string{
		r.Header.Get("X-Client-Real-Ip"),
		r.Header.Get("X-Real-Ip"),
		r.Header.Get("CF-Connecting-IP"),
		r.Header.Get("True-Client-Ip"),
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			lastIP := strings.TrimSpace(ips[len(ips)-1])
			if lastIP != "" {
				return lastIP
			}
		}
		//for _, ip := range strings.Split(xff, ",") {
		//	ip = strings.TrimSpace(ip)
		//	if ip != "" && !isPrivateIP(net.ParseIP(ip)) {
		//		return ip
		//	}
		//}
	}

	for _, ip := range ipSources {
		if ip != "" && !isPrivateIP(net.ParseIP(ip)) {
			return ip
		}
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func isPrivateIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	if ip4 := ip.To4(); ip4 != nil {
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	return false
}
