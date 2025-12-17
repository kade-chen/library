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
	}).
		Doc("Swagger UI").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	// ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
	// 	//2.restfulspec.BuildSwagger() 方法使用这个配置来生成对应的 Swagger 文档
	// 	swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
	// 	w.WriteAsJson(swagger)
	// }))

	ws.Route(ws.GET("/ui").To(h.SwaggerUI).
		Doc("Swagger reddoc UI").
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
	if swagger == nil {
		return
	}

	// 1. 先删掉 definitions 里的非法 structpb 定义
	if swagger.Definitions != nil {
		delete(swagger.Definitions, "structpb.isValue_Kind")
		delete(swagger.Definitions, "structpb.Value")
		delete(swagger.Definitions, "structpb.Struct")
	}

	// 2. 扫描所有 schema，把指向 structpb 的 $ref 改成 object
	// fixSchema := func(s *spec.Schema) {}
	var fixSchemaFn func(s *spec.Schema)

	fixSchemaFn = func(s *spec.Schema) {
		if s == nil {
			return
		}

		// 命中非法 $ref
		ref := s.Ref.String()
		if ref == "#/definitions/structpb.Struct" ||
			ref == "#/definitions/structpb.Value" ||
			ref == "#/definitions/structpb.isValue_Kind" {

			// 🔥 直接整体替换成 object
			*s = spec.Schema{
				SchemaProps: spec.SchemaProps{
					Type: []string{"object"},
					AdditionalProperties: &spec.SchemaOrBool{
						Allows: true,
					},
				},
			}
			return
		}

		// properties
		for k := range s.Properties {
			prop := s.Properties[k]
			fixSchemaFn(&prop)
			s.Properties[k] = prop
		}

		// array items
		if s.Items != nil && s.Items.Schema != nil {
			fixSchemaFn(s.Items.Schema)
		}

		// allOf / anyOf / oneOf
		for i := range s.AllOf {
			fixSchemaFn(&s.AllOf[i])
		}
		for i := range s.AnyOf {
			fixSchemaFn(&s.AnyOf[i])
		}
		for i := range s.OneOf {
			fixSchemaFn(&s.OneOf[i])
		}

		// additionalProperties
		if s.AdditionalProperties != nil && s.AdditionalProperties.Schema != nil {
			fixSchemaFn(s.AdditionalProperties.Schema)
		}
	}

	// 3. definitions 本身也要递归修
	for name, def := range swagger.Definitions {
		fixSchemaFn(&def)
		swagger.Definitions[name] = def
	}

	// 4. paths / operations / parameters / responses 全量扫描
	for _, path := range swagger.Paths.Paths {
		for _, op := range []*spec.Operation{
			path.Get,
			path.Post,
			path.Put,
			path.Delete,
			path.Patch,
			path.Options,
			path.Head,
		} {
			if op == nil {
				continue
			}

			// parameters
			for i := range op.Parameters {
				if op.Parameters[i].Schema != nil {
					fixSchemaFn(op.Parameters[i].Schema)
				}
			}

			// responses
			for code, resp := range op.Responses.StatusCodeResponses {
				if resp.Schema != nil {
					fixSchemaFn(resp.Schema)
				}
				op.Responses.StatusCodeResponses[code] = resp
			}
		}
	}
}
