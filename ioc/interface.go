package ioc

import "context"

type Stroe interface {
	StoreUser
	StoreManage //暂时不需要环境变量加载配置
}

type StoreUser interface {
	// 对象注册
	Registry(obj Object)
	// 对象获取
	Get(name string, opts ...GetOption) Object
	// 根据对象类型, 直接加载对象
	Load(obj any, opts ...GetOption) (Object, error) //这个接口有问题不完善，后续需要优化
	// 打印对象列表
	List() []string
	// 数量统计
	Count() int
	// 遍历注入的对象 这个是给加载配置文件用的
	ForEach(fn func(*ObjectWrapper))
}

type GetOption func(*option)

type option struct {
	version string
}

// 这个可以优化一下，并不能实现版本
func (o *option) Apply(opts ...GetOption) *option {
	for i := range opts {
		// fmt.Println("-------", opts)
		opts[i](o)
	}
	// fmt.Println("-------", o)
	return o
}

func defaultOption() *option {
	return &option{
		version: DEFAULT_VERSION,
	}
}

type StoreManage interface {
	// 从环境变量中加载对象配置
	LoadFromEnv(prefix string) error
}

// Object 内部服务实例, 不需要暴露
type Object interface {
	// 对象初始化
	Init() error
	// 对象的名称
	Name() string
	// 对象版本
	Version() string
	// 对象优先级
	Priority() int
	// 对象的销毁
	Close(ctx context.Context) error
	// 是否允许同名对象被替换, 默认不允许被替换.根据注册先后
	AllowOverwrite() bool
	// // 对象一些元数据, 对象的更多描述信息, 扩展使用
	Meta() ObjectMeta
}
