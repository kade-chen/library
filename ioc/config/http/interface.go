package http

import "github.com/kade-chen/library/ioc"

const (
	AppName = "http"
)

func Get() *Http {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Http)
}
