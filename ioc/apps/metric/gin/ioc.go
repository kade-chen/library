package metric

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/apps/metric"
	"github.com/kade-chen/library/ioc/config/log"

	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&ginHandler{
		Metric: metric.NewDefaultMetric(),
	})
}

type ginHandler struct {
	log *zerolog.Logger
	ioc.ObjectImpl

	*metric.Metric
}

func (h *ginHandler) Init() error {
	h.log = log.Sub(metric.AppName)
	h.Registry()
	return nil
}

func (h *ginHandler) Name() string {
	return metric.AppName
}

func (h *ginHandler) Version() string {
	return "v1"
}

func (h *ginHandler) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Endpoint
	return meta
}
