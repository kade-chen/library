package gorestful

import (
	"github.com/kade-chen/library/ioc"
	"github.com/emicklei/go-restful/v3"
	"github.com/rs/zerolog"

	"github.com/kade-chen/library/ioc/config/application"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/kade-chen/library/ioc/config/trace"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func init() {
	ioc.Config().Registry(&GoRestfulFramework{})
}

type GoRestfulFramework struct {
	ioc.ObjectImpl
	log       *zerolog.Logger
	container *restful.Container
}

func (i *GoRestfulFramework) Name() string {
	return AppName
}

func (i *GoRestfulFramework) Priority() int {
	return 9996
}

func (i *GoRestfulFramework) Init() error {
	i.log = log.Sub(AppName)
	i.container = restful.DefaultContainer

	// 设置默认的content-type为json
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	if http.Get().EnableTrace && trace.Get().Enable {
		i.log.Info().Msg("enable go-restful trace")
		i.container.Filter(otelrestful.OTelFilter(application.Get().GetAppNameWithDefault("default")))
		// Filter  添加一个过滤器的方法。过滤器用于在请求处理过程中执行一些额外的逻辑，如日志记录、身份验证、监控等。
		// otelrestful.OTelFilter() 用于创建一个 OpenTelemetry 过滤器。OpenTelemetry 是一个用于服务端观测的开源框架，可以收集和报告指标、日志和分布式追踪数据。
	}

	// 把container路由器 注册给Http服务器
	http.Get().SetRouter(i.container)

	return nil
}
