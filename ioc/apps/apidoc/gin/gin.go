package gin

import (
	"github.com/kade-chen/library/ioc/config/gogin"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

func (h *SwaggerApiDoc) Registry() {
	r := gogin.InitRouter(h)
	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteString(swag.GetSwagger(h.InstanceName).ReadDoc())
	})

	h.log.Info().Msgf("Get the API Doc using http://%s", http.Get().ApiObjectAddr(h))
}
