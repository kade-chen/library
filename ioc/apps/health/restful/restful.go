package restful

import (
	"github.com/kade-chen/library/ioc/config/gorestful"
	"github.com/kade-chen/library/ioc/config/http"

	"github.com/kade-chen/library/http/response"
	ioc_health "github.com/kade-chen/library/ioc/apps/health"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
)

func (h *HealthChecker) Registry() {

	ws := gorestful.InitRouter(h)
	tags := []string{"健康检查"}
	ws.Route(ws.GET("/").To(h.HealthHandleFunc).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", ioc_health.HealthCheckResponse{}))

	h.log.Info().Msgf("Get the Health using http://%s%s", http.Get().Addr(), h.Path)
}

func (h *HealthChecker) HealthHandleFunc(r *restful.Request, w *restful.Response) {
	req := ioc_health.NewHealthCheckRequest()
	resp, err := h.Service.Check(
		r.Request.Context(),
		req,
	)
	if err != nil {
		response.Failed(w, err)
		return
	}
	//如果健康检查成功，则将响应转换为 JSON 格式，并通过 w 写入到响应中。
	err = w.WriteAsJson(ioc_health.NewHealth(resp))
	if err != nil {
		h.log.Error().Msgf("send success response error, %s", err)
	}
}
