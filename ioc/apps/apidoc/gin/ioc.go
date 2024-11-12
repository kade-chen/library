package gin

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/apps/apidoc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.ApiDoc{
			Path: "/swagger.json",
		},
		InstanceName: "swagger",
	})
}

type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	apidoc.ApiDoc
	InstanceName string `json:"instance_name" yaml:"instance_name" toml:"instance_name" env:"SWAGGER_INSTANCE_NAME"`
}

func (h *SwaggerApiDoc) Name() string {
	return apidoc.AppNamegin
}

func (h *SwaggerApiDoc) Init() error {
	h.log = log.Sub("api_doc")
	h.Registry()
	return nil
}

func (i *SwaggerApiDoc) Priority() int {
	return -100
}

func (h *SwaggerApiDoc) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Path
	return meta
}
