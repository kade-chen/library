package ioc

const (
	DEFAULT_NAMESPACE = "default"
)

// 默认空间, 用于托管工具类, 在控制器之前进行初始化
// store.Namespace(DEFAULT_NAMESPACE) 这个函数的返回，*NamespaceStore实现了 StoreUser 接口，所以可以作为 StoreUser 类型返回。
func Default() StoreUser {
	return store.Namespace(DEFAULT_NAMESPACE)
}
