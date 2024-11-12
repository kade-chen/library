package ioc

// ioc服务注册需要的结构体
// string 是ioc的名字
// any是服务的接口
type iocContainer struct {
	//采片Map 来保持对象注冊
	store map[string]iocObject
}


