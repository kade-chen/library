之前的ioc
map形式，注册name：“interface"
如果掉用其方法，就是实现init，这样才调

现在的ioc
切片，只需要注册 name，version
如果掉用其方法，就是实现init，这样才调，但是这种方式只要name和version，比上面更复杂好用


1.使用load前要注册
ioc.Default().Registry(&c{})
ioc.Default().Load()
