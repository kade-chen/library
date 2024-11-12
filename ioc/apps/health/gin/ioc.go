package gin

import (
	"github.com/kade-chen/library/ioc"
	ioc_health "github.com/kade-chen/library/ioc/apps/health"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/health"
	healthgrpc "google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	ioc.Api().Registry(&HealthChecker{
		HealthCheck: ioc_health.HealthCheck{
			Path: ioc_health.DEFAUL_HEALTH_PATH,
		},
	})
}

type HealthChecker struct {
	ioc.ObjectImpl
	ioc_health.HealthCheck
	log *zerolog.Logger

	Service healthgrpc.HealthServer `ioc:"autowire=true;namespace=controllers"`
}

func (h *HealthChecker) Name() string {
	return ioc_health.AppName
}

func (h *HealthChecker) Init() error {
	if h.Service == nil {
		h.Service = health.NewServer()
	}
	h.log = log.Sub("health_check")
	h.Registry()
	return nil
}

func (h *HealthChecker) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Path
	return meta
}
