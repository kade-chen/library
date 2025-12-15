package ip

import (
	"net"
	"net/http"
	"strings"
)

var (
	// DefaultScanForwareHeaderKey 协商forward ip 的 hander key名称
	DefaultScanForwareHeaderKey = []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP", // Cloudflare
		"X-Client-IP",
		"X-Original-Forwarded-For",
	}
)

// GetRemoteIP todo
func GetRemoteIP(r *http.Request) (proxyIP, remoteIP string) {
	// 1️⃣ 解析 RemoteAddr（真实 TCP 对端 IP）
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		remoteIP = host
	} else {
		remoteIP = r.RemoteAddr
	}

	// 2️⃣ 从代理头里取客户端 IP
	for _, key := range DefaultScanForwareHeaderKey {
		value := r.Header.Get(key)
		if value == "" {
			continue
		}

		// X-Forwarded-For: client, proxy1, proxy2
		parts := strings.Split(value, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			proxyIP = ip
			break
		}
	}

	return
}
