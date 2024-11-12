package metric

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/apps/metric"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&restfulHandler{
		Metric: metric.NewDefaultMetric(),
	})
}

type restfulHandler struct {
	log *zerolog.Logger
	ioc.ObjectImpl

	*metric.Metric
}

func (h *restfulHandler) Init() error {
	h.log = log.Sub(metric.AppName)
	h.Registry()
	return nil
}

func (h *restfulHandler) Name() string {
	return metric.AppName
}

func (h *restfulHandler) Version() string {
	return "v1"
}

func (h *restfulHandler) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Endpoint
	return meta
}
