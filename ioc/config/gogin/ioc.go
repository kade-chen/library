package gogin

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/application"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/kade-chen/library/ioc/config/trace"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func init() {
	ioc.Config().Registry(&GinFramework{})
}

type GinFramework struct {
	ioc.ObjectImpl
	Engine *gin.Engine
	log    *zerolog.Logger
}

func (g *GinFramework) Init() error {
	g.log = log.Sub(g.Name())
	g.Engine = gin.Default()

	if http.Get().EnableTrace && trace.Get().Enable {
		g.log.Info().Msg("enable gin trace")
		g.Engine.Use(otelgin.Middleware(application.Get().GetAppNameWithDefault("default")))
	}

	// 注册给Http服务器
	http.Get().SetRouter(g.Engine)
	return nil
}

func (g *GinFramework) Priority() int {
	return 9996
}

func (g *GinFramework) Name() string {
	return AppName
}
