package httpgorestful

import (
	"context"
	"testing"

	"gitee.com/go-kade/library/ioc"
	_ "gitee.com/go-kade/library/ioc/apps/apidoc/restful"
	_ "gitee.com/go-kade/library/ioc/apps/health/restful"
	_ "gitee.com/go-kade/library/ioc/apps/metric/restful"
	_ "gitee.com/go-kade/library/ioc/config/cors"
	"gitee.com/go-kade/library/ioc/server"

	"github.com/emicklei/go-restful/v3"
)

func TestDefaultConfig(t *testing.T) {
	// 注册HTTP接口类
	ioc.Api().Registry(&HelloServiceApiHandler{})
	// 启动应用
	err := server.Run(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

type HelloServiceApiHandler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl
}

// 模块的名称, 会作为路径的一部分比如: /mcube_service/api/v1/hello_module/
// 路径构成规则 <service_name>/<path_prefix>/<service_version>/<module_name>
// http://127.0.0.1:8080/default/api/v1/hello_module
func (h *HelloServiceApiHandler) Name() string {
	return "hello_module"
}

func (h *HelloServiceApiHandler) Version() string {
	return "v1"
}

// API路由
func (h *HelloServiceApiHandler) Registry(ws *restful.WebService) {
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(map[string]string{
			"data": "hello mcube1111",
		})
	}))
}
