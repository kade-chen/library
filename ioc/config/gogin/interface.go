package gogin

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/http"
	"github.com/gin-gonic/gin"
)

const (
	AppName = "gin_webframework"
)

func RootRouter() *gin.Engine {
	return ioc.Config().Get(AppName).(*GinFramework).Engine
}

// 添加对象前缀路径
func InitRouter(obj ioc.Object) gin.IRouter {
	modulePath := http.Get().ApiObjectPathPrefix(obj)
	return ioc.Config().Get(AppName).(*GinFramework).Engine.Group(modulePath)
}
