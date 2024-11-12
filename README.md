#### Links
- [Setup one](https://pkg.go.dev/github.com/kade-chen/library/v1/tree/v1.0.2)
- [Setup two](https://pkg.go.dev/github.com/kade-chen/library/tree/v1/v1.0.2)
- [Setup three](https://pkg.go.dev/github.com/kade-chen/library/v1@v1.0.2)
- [Github Versions](https://pkg.go.dev/github.com/kade-chen/library?tab=versions)

# 之前的ioc
map形式，注册name：“interface"

如果掉用其方法，就是实现init，这样才调

# 现在的ioc
切片，只需要注册 name，version

如果掉用其方法，就是实现init，这样才调，但是这种方式只要name和version，比上面更复杂好用

# ioc 2.0
优化隐试接口

### app1 ioc1 是老的版本  map

### app  ioc  是新的版本  slice

#### 补充业务异常
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "/Users/kade.chen/go12-project/sso/etc/application.toml"
	ioc.DevelopmentSetup(req)
自动加载注册init所有ioc


# 更快到官网
```go
 git tag v1.0.4

 git push -f https://github.com/kade-chen/library.git v1.0.2  cobra command successful

 https://pkg.go.dev/github.com/kade-chen/library/v1/tree/v1.0.2

//  https://pkg.go.dev/github.com/kade-chen/library/tree/v1/v1.0.2

 https://pkg.go.dev/github.com/kade-chen/library/v1@v1.0.2 目前来说，v2版本只需要这个就行了，前面两个不需要

备注：https://deps.dev/search?q=github.com%2Fkade-chen&system=go
如果不知道版本可以从这里看


最新 13分钟同步完成  21分钟前 https://pkg.go.dev/github.com/kade-chen/library?tab=versions 这个出来

 
```

```go
 git tag v2.1.3

 git push -f https://github.com/kade-chen/library.git v2.1.3  cobra command successful

 https://pkg.go.dev/github.com/kade-chen/library/v2/tree/v2.1.3

//  https://pkg.go.dev/github.com/kade-chen/library/tree/v2/v2.1.3

 https://pkg.go.dev/github.com/kade-chen/library/v2@v2.1.3 目前来说，v2版本只需要这个就行了，前面两个不需要

备注：https://deps.dev/go/gitee.com%2Fgo-kade%2Flibrary/v1.0.1-0.20240201092113-6e4b7db5c891/versions
如果不知道版本可以从这里看


最新 13分钟同步完成  21分钟前 https://deps.dev/go/gitee.com%2Fgo-kade%2Flibrary/v1.2.1/versions 这个出来

 
```


```go

v1 before

git tag v1.0.0

git push -f https://github.com/kade-chen/library.git v1.0.0

https://pkg.go.dev/github.com/kade-chen/library/tree/v1.0.0

https://pkg.go.dev/github.com/kade-chen/library@v1.0.0

更新大版本，例如 v1.0.0 -> v2.0.0


 go mod edit -module  github.com/kade-chen/library/v2

find . -type f -name '*.go' \
    -exec sed -i '' -e 's,github.com/kade-chen/library,github.com/kade-chen/library/v2,g' {} \;


 现在我们有了一个v2模块，但我们想在发布版本之前进行试验和更改。在我们发布v2.0.0（或任何没有预发布后缀的版本）之前，我们可以在决定新 API 时开发和进行重大更改。如果我们希望用户能够在我们正式稳定之前试用新 API，我们可以发布v2预发布版本：
 git tag v2.0.0-alpha.1
 git push https://github.com/kade-chen/library.git v2.0.0-alpha.1

一旦我们对我们的 API 感到满意v2，并且确定我们不需要任何其他重大更改，我们就可以标记v2.0.0：
 git tag v2.0.0

 git push -f https://github.com/kade-chen/library.git v2.0.0

 https://pkg.go.dev/github.com/kade-chen/library/v2/tree/v2.0.0

 https://pkg.go.dev/github.com/kade-chen/library/v2@v2.0.0
 
```


# 
_ "github.com/kade-chen/library/ioc/config/cors" 

_ "github.com/kade-chen/library/ioc/apps/apidoc/restful"

    想用apidoc得闲注入cors

	_ "github.com/kade-chen/library/ioc/apps/apidoc/restful" //依赖cors，不然访问会失败，里面只有apidoc
	_ "github.com/kade-chen/library/ioc/apps/health/restful"  
    _ "github.com/kade-chen/library/ioc/apps/metric/restful"  //promethus的指标暴露
	_ "github.com/kade-chen/library/ioc/config/cors"  //如果不加这个，跨域可能会有问题，另外apidoc也用不了

    _ "github.com/kade-chen/library/ioc/apps/apidoc/swaggo" gin

```go

go-restful

	// 开启apidoc 必须开启cors
	_ "github.com/kade-chen/library/ioc/apps/apidoc/restful"

	// 开启Health健康检查
	_ "github.com/kade-chen/library/ioc/apps/health/restful"

	// 开启Metric
	// promethus的指标暴露
	_ "github.com/kade-chen/library/ioc/apps/metric/restful"

	// 开启CORS, 允许资源跨域共享
	_ "github.com/kade-chen/library/ioc/config/cors/gorestful"

```


```go

go-gin

	// 引入生成好的API Doc代码 必须开启cors
	_ "github.com/kade-chen/library/examples/http_gin/docs"

	// 引入集成工程
	_ "github.com/kade-chen/library/ioc/apps/apidoc/gin"

	// 开启Health健康检查
	_ "github.com/kade-chen/library/ioc/apps/health/gin"

	// 开启Metric
	// promethus的指标暴露
	_ "github.com/kade-chen/library/ioc/apps/metric/gin"

	// 开启CORS, 允许资源跨域共享
	_ "github.com/kade-chen/library/ioc/config/cors/gin"

```


```go

package main

import (
	"context"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/server"
	"github.com/emicklei/go-restful/v3"

	// 开启apidoc
	_ "github.com/kade-chen/library/ioc/apps/apidoc/restful"

	// 开启Health健康检查
	_ "github.com/kade-chen/library/ioc/apps/health/restful"
	// 开启Metric
	_ "github.com/kade-chen/library/ioc/apps/metric/restful"
	// 开启CORS, 允许资源跨域共享
	_ "github.com/kade-chen/library/ioc/config/cors/gorestful"

	// _ "github.com/kade-chen/library/ioc/apps/apidoc/swaggo" gin
	_ "github.com/kade-chen/library/ioc/config/cors/gorestful"
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
	h.Registry()
	return nil
}

// API路由
func (h *HelloServiceApiHandler) Registry() error  {
	ws := gorestful.InitRouter(h)
	///default/api/v1/hello_module/
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(map[string]string{
			"data": "hello mcube1111111",
		})
	}))
	//default/api/v1/hello_module/cc
	ws.Route(ws.GET("/cc").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(map[string]string{
			"data": "hello mcube1111111",
		})
	}))
	//default/api/v1/hello_module/cc/111
	ws.Route(ws.GET("/cc/111").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(map[string]string{
			"data": "hello mcube1111111",
		})
	}))
}

//自定义路由
func (i *SwaggerApiDoc) Meta() ioc.ObjectMeta {
	return DefaultObjectMeta1()
}

func DefaultObjectMeta1() ioc.ObjectMeta {
	return ioc.ObjectMeta{
		//		CustomPathPrefix: "/s", 必须要/号 http://127.0.0.1:8080/s
		CustomPathPrefix: "/swagger.json",
		// CustomPathPrefix: "/s",
		Extra: map[string]string{},
	}
}

```


```go
package main

import (
	"context"
	"net/http"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/server"
	"github.com/gin-gonic/gin"

	// 引入生成好的API Doc代码
	_ "github.com/kade-chen/library/examples/http_gin/docs"
	// 引入集成工程
	_ "github.com/kade-chen/library/ioc/apps/apidoc/gin"
	// 开启Health健康检查
	_ "github.com/kade-chen/library/ioc/apps/health/gin"
	// 开启Metric
	_ "github.com/kade-chen/library/ioc/apps/metric/gin"
	// 开启CORS, 允许资源跨域共享
	_ "github.com/kade-chen/library/ioc/config/cors/gin"
	"github.com/kade-chen/library/ioc/config/gogin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
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
func (h *HelloServiceApiHandler) Name() string {
	return "hello_module"
}

func (h *HelloServiceApiHandler) Version() string {
	return "v1"
}

func (h *HelloServiceApiHandler) Init() error {
	h.Registry()
	return nil
}

// API路由
func (h *HelloServiceApiHandler) Registry() {

	r := gogin.InitRouter(h)
	///default/api/v1/hello_module/
	r.GET("/", h.Hello)
	//default/api/v1/hello_module/cc
	r.GET("/cc", h.Hello)
	///default/api/v1/hello_module/cc/111
	r.GET("/cc/111", h.Hello)
}

// @Summary 修改文章标签
// @Description  修改文章标签
// @Tags         文章管理
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func (h *HelloServiceApiHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "hello mcube",
	})
}

func (i *HelloServiceApiHandler) Meta() ioc.ObjectMeta {
	return ioc.ObjectMeta(DefaultObjectMeta())
}

func DefaultObjectMeta() ObjectMeta {
	return ObjectMeta{
		//		CustomPathPrefix: "/s", 必须要/号 http://127.0.0.1:8080/s
		CustomPathPrefix: "/bb",
		// CustomPathPrefix: "/s",
		Extra: map[string]string{},
	}
}

type ObjectMeta struct {
	CustomPathPrefix string
	Extra            map[string]string
}

```