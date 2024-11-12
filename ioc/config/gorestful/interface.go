package gorestful

import (
	"fmt"

	"github.com/kade-chen/library/http/restful/accessor/form"
	"github.com/kade-chen/library/http/restful/accessor/yaml"
	"github.com/kade-chen/library/http/restful/accessor/yamlk8s"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/emicklei/go-restful/v3"
)

const (
	AppName = "restful_webframework"
)

func RootRouter() *restful.Container {
	return ioc.Config().Get(AppName).(*GoRestfulFramework).container
}

func InitRouter(obj ioc.Object) *restful.WebService {
	ws := new(restful.WebService)
	ws.
		// Path(ApiPathPrefix(pathPrefix, api.Meta(), api))
		//path 根路径，api.Registry(ws)中的get是子路径
		Path(ApiPathPrefix(obj)).
		Consumes(restful.MIME_JSON, form.MIME_POST_FORM, form.MIME_MULTIPART_FORM, yaml.MIME_YAML, yamlk8s.MIME_YAML).
		Produces(restful.MIME_JSON, yaml.MIME_YAML, yamlk8s.MIME_YAML)
	// 添加到Root Container
	RootRouter().Add(ws)
	return ws
}

// func ApiPathPrefix(pathPrefix string, obm ObjectMeta, obj Object) string {
func ApiPathPrefix(obj ioc.Object) string {
	fmt.Println("-----path:", http.Get().Addr() + http.Get().ApiObjectPathPrefix(obj))
	return http.Get().ApiObjectPathPrefix(obj)
}
