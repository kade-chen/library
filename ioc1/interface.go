package ioc

// 定义注册传进来的app_name
// 定义注册进来的接口以及方法，初始化其实就是把对象注册进去
type iocObject interface {
	//注册对象名称
	Name() string
	//注册对象初始化
	ServiceInit() error
}

// ioc服务注册,obj.Name()名字等于接口
func (i *iocContainer) Registry(obj iocObject) {
	i.store[obj.Name()] = obj
}

// 获取ioc对应服务的接口
func (i *iocContainer) GetName(name string) any {
	return i.store[name]
}

// 负责初始化所有对象
func (i *iocContainer) Init() error {
	for _, v := range i.store {
		if err := v.ServiceInit(); err != nil {
			return err
		}
	}
	return nil
}
