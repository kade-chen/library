package metric

import (
	"github.com/kade-chen/library/ioc/config/gorestful"
	"github.com/kade-chen/library/ioc/config/http"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (h *restfulHandler) Registry() {
	tags := []string{"指标"}

	ws := gorestful.InitRouter(h)

	ws.Route(ws.GET("/").
		To(h.MetricHandleFunc).
		Doc("创建Job").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	h.log.Info().Msgf("Get the Metric using http://%s%s", http.Get().Addr(), h.Endpoint)
}

func (h *restfulHandler) MetricHandleFunc(r *restful.Request, w *restful.Response) {
	// 基于标准库 包装了一层
	promhttp.Handler().ServeHTTP(w, r.Request)
}
