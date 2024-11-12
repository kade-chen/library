package main

import (
	"context"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/server"
	"github.com/emicklei/go-restful/v3"

	_ "github.com/kade-chen/library/ioc/apps/apidoc/restful"

	// 开启Health健康检查
	// _ "github.com/kade-chen/library/ioc/apps/health/restful"
	// 开启Metric
	// _ "github.com/kade-chen/library/ioc/apps/metric/restful"
	// 开启CORS, 允许资源跨域共享
	_ "github.com/kade-chen/library/ioc/config/cors/gorestful"

	// _ "github.com/kade-chen/library/ioc/apps/apidoc/swaggo" gin
	// _ "github.com/kade-chen/library/ioc/config/cors/gorestful"
	"github.com/kade-chen/library/ioc/config/gorestful"
)

func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&HelloServiceApiHandler{})
	// 启动应用
	// req := ioc.NewLoadConfigRequest()
	// req.ConfigFile.Enabled = true
	// req.ConfigFile.Path = "/Users/kade.chen/go12-project/library/etc/application.toml"
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

func (h *HelloServiceApiHandler) Init() error {
	h.Registry2()
	return nil
}

// API路由
// func (h *HelloServiceApiHandler) Registry1(ws *restful.WebService) {
// 	///default/api/v1/hello_module/
// 	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
// 		w.WriteAsJson(map[string]string{
// 			"data": "hello mcube1111111",
// 		})
// 	}))
// 	//default/api/v1/hello_module/cc
// 	ws.Route(ws.GET("/cc").To(func(r *restful.Request, w *restful.Response) {
// 		w.WriteAsJson(map[string]string{
// 			"data": "hello mcube1111111",
// 		})
// 	}))
// 	//default/api/v1/hello_module/cc/111
// 	ws.Route(ws.GET("/cc/111").To(func(r *restful.Request, w *restful.Response) {
// 		w.WriteAsJson(map[string]string{
// 			"data": "hello mcube1111111",
// 		})
// 	}))
// }

func (h *HelloServiceApiHandler) Registry2() error {
	ws := gorestful.InitRouter(h)
	///default/api/v1/hello_module/
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(map[string]string{
			"data": "hello mcube1111111",
		})
	}))
	return nil
}
