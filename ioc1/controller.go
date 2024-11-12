package ioc


// 专门用于注册Controller对象
func Controller() *iocContainer {
	return iocImpl
}
// ioc 注册表对象, 全局只有
var iocImpl = &iocContainer{
	store: map[string]iocObject{},
}

// 专门用于注册api-Controller对象
func ApiHandler() *iocContainer {
	return apiHandlerContainer
}

// api-ioc 注册表对象, 全局只有
var apiHandlerContainer = &iocContainer{
	store: map[string]iocObject{},
}