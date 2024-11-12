package gin

import (
	"github.com/kade-chen/library/ioc/config/gogin"
	"github.com/kade-chen/library/ioc/config/http"

	h_response "github.com/kade-chen/library/http/response"
	ioc_health "github.com/kade-chen/library/ioc/apps/health"
	"github.com/gin-gonic/gin"
)

func (h *HealthChecker) Registry() {

	r := gogin.InitRouter(h)

	r.GET("/", h.HealthHandleFunc)

	h.log.Info().Msgf("Get the Health using http://%s%s", http.Get().Addr(), h.Path)

}



// @Summary 健康检查
// @Description  修改文章标签
// @Tags         文章管理
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /heathcheck [put]
func (h *HealthChecker) HealthHandleFunc(c *gin.Context) {
	req := ioc_health.NewHealthCheckRequest()
	resp, err := h.Service.Check(
		c.Request.Context(),
		req,
	)
	if err != nil {
		h_response.Failed(c.Writer, err)
		return
	}

	h_response.Success(c.Writer, ioc_health.NewHealth(resp))
}
