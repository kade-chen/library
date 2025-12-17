package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/apps/apidoc"
	"github.com/kade-chen/library/ioc/config/gorestful"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"

	// 开启apidoc 必须开启cors
	_ "github.com/kade-chen/library/ioc/config/cors/gorestful"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.ApiDoc{
			Path: "/swagger.json",
		},
	})
}

type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	apidoc.ApiDoc
	// Path string `json:"path" yaml:"path" toml:"path" env:"HTTP_API_DOC_PATH"`
}

func (h *SwaggerApiDoc) Name() string {
	return apidoc.AppName
}

func (h *SwaggerApiDoc) Init() error {
	h.log = log.Sub("api_doc")
	h.Registry()
	return nil
}

func (i *SwaggerApiDoc) Priority() int {
	return -100
}

func (i *SwaggerApiDoc) Meta() ioc.ObjectMeta {
	ObjectMeta := ioc.DefaultObjectMeta()
	ObjectMeta.CustomPathPrefix = i.Path
	return ObjectMeta
}

func (h *SwaggerApiDoc) Registry() {
	tags := []string{"API 文档"}
	ws := gorestful.InitRouter(h)

	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())

		// 🔥 关键一步：patch
		patchSwagger(swagger)

		w.WriteAsJson(swagger)
	}))
	// ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
	// 	//2.restfulspec.BuildSwagger() 方法使用这个配置来生成对应的 Swagger 文档
	// 	swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
	// 	w.WriteAsJson(swagger)
	// }))

	ws.Route(ws.GET("/ui").To(h.SwaggerUI).
		Doc("Swagger UI").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	// h.log.Info().Msgf("Get the API UI using %s", h.ApiUIPath())

	if h.Meta().CustomPathPrefix != "" {
		h.log.Info().Msgf("Get the API Doc using http://%s%s", http.Get().Addr(), http.Get().ApiObjectPathPrefix(h))
	} else {
		h.log.Info().Msgf("Get the API Doc using http://%s%s", http.Get().Addr(), http.Get().ApiObjectAddr(h))
	}
}

func patchSwagger(swagger *spec.Swagger) {
	if swagger.Definitions == nil {
		return
	}

	// 1. 删除 structpb 相关的非法定义
	delete(swagger.Definitions, "structpb.isValue_Kind")
	delete(swagger.Definitions, "structpb.Value")
	delete(swagger.Definitions, "structpb.Struct")

	// 2. 可选：兜底，把引用它们的地方统一改成 object
	for _, def := range swagger.Definitions {
		if def.SchemaProps.Type.Contains("object") {
			continue
		}
	}
}
