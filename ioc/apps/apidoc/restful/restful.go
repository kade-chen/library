package restful

import (
	"fmt"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/kade-chen/library/ioc/apps/apidoc"
	"github.com/kade-chen/library/ioc/config/http"
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

func (h *SwaggerApiDoc) SwaggerUI(r *restful.Request, w *restful.Response) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(apidoc.HTML_REDOC, "kade-path")))
}

// func (h *SwaggerApiDoc) ApiDocPath() string {
// 	if application.Get().AppAddress != "" {
// 		return application.Get().AppAddress + filepath.Join(http.Get().ApiObjectPathPrefix(h), h.JsonPath)
// 	}

// 	return http.Get().ApiObjectAddr(h) + h.JsonPath
// }
