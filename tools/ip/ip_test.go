package ip_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/kade-chen/google-billing-console/tools/ip"

	"github.com/stretchr/testify/assert"
)

func TestGetRemoteIpFromHeader(t *testing.T) {
	shoud := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "10.10.10.10")
	ip1, ip := ip.GetRemoteIP(req)
	shoud.Equal("10.10.10.10", ip1, ip)

	fmt.Println(ip, shoud.Equal("10.10.10.10", ip))
}

func TestGetRemoteIpFromConn(t *testing.T) {
	shoud := assert.New(t)

	req, _ := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.10.10.10"
	ip1, ip := ip.GetRemoteIP(req)
	shoud.Equal("10.10.10.10", ip, ip1)

	fmt.Println(ip, shoud.Equal("10.10.10.10", ip))
}
