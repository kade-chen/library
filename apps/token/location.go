package token

import (
	"net/http"

	tools "github.com/kade-chen/google-billing-console/tools/ip"
	"github.com/mssola/user_agent"
)

func NewLocation() *Location {
	return &Location{
		IpLocation: &IPLocation{},
		UserAgent:  &UserAgent{},
	}
}

func NewNewLocationFromHttp(r *http.Request) *Location {
	l := NewLocation()

	// 解析UserAgent
	ua := r.UserAgent()
	if ua != "" {
		ua := user_agent.New(ua)
		l.UserAgent = &UserAgent{
			Os:       ua.OS(),
			Platform: ua.Platform(),
		}
		l.UserAgent.EngineName, l.UserAgent.EngineVersion = ua.Engine()
		l.UserAgent.BrowserName, l.UserAgent.BrowserVersion = ua.Browser()
	}

	// 解析地理位置
	l.IpLocation.ProxyIp, l.IpLocation.RemoteIp = tools.GetRemoteIP(r)
	return l
}
