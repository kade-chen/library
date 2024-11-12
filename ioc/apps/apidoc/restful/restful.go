package restful

import (
	"github.com/kade-chen/library/ioc/config/http"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

// API Doc config
func (h *SwaggerApiDoc) SwaggerDocConfig() restfulspec.Config {
	return restfulspec.Config{
		// 1.是用于获取当前已注册的 Web 服务列表的函数
		WebServices: restful.RegisteredWebServices(),
		APIPath:     http.Get().ApiObjectPathPrefix(h),
		// APIPath:                       http.Get().ApiObjectPathPrefix(h),
		PostBuildSwaggerObjectHandler: http.Get().SwagerDocs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}
}
