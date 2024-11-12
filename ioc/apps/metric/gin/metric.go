package metric

import (
	"github.com/kade-chen/library/ioc/config/gogin"
	"github.com/kade-chen/library/ioc/config/http"

	"github.com/gin-gonic/gin"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *ginHandler) Registry() {

	r := gogin.InitRouter(h)

	r.GET("/", func(ctx *gin.Context) {
		// 基于标准库 包装了一层
		promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
	})

	h.log.Info().Msgf("Get the Metric using http://%s%s", http.Get().Addr(), h.Endpoint)
}
