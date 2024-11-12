package ioc

import "github.com/gin-gonic/gin"

// gin http的ioc路由注册
type GinApiHandler interface {
	Registry(r gin.IRouter)
}

// 管理者所有的对象(Api Handler)
// 把每个 ApiHandler的路由注册给Root Router
func (c *iocContainer) RouteRegistry(r gin.IRouter) {
	// 找到被托管的APIHandler
	for _, obj := range c.store {
		// 怎么来判断这个对象是一个APIHandler对象
		if api, ok := obj.(GinApiHandler); ok { //断言这个对象有没有实现（Registry(r gin.IRouter)） 这个接口的这个方法，如果实现了这个接口，把它注册进去
			api.Registry(r)
		}
	}
}

