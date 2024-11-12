package ioc

import (
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type ObjectWrapper struct {
	Name           string
	Version        string
	AllowOverwrite bool
	Priority       int
	Value          Object
}

func NewObjectWrapper(obj Object) *ObjectWrapper {
	name, version := GetIocObjectUid(obj)
	//在load里面做了
	// err := obj.Init()
	// if err != nil {
	// 	panic(err)
	// }
	return &ObjectWrapper{
		Name:           name,
		Version:        version,
		Priority:       obj.Priority(),
		AllowOverwrite: obj.AllowOverwrite(),
		Value:          obj,
	}
}

var _ StoreUser = (*NamespaceStore)(nil)
var _ Stroe = (*NamespaceStore)(nil)
// var _ Object = (*NamespaceStore)(nil)

// 这个初始化，在store的时候，也就是在namespace.go中
type NamespaceStore struct {
	// 空间名称
	Namespace string
	// 空间优先级
	Priority int
	// 空间对象列表
	Items []*ObjectWrapper
}

// -----按照priority排序注册
func (s *NamespaceStore) Sort() {
	sort.Sort(s)
}

// Len 返回切片的长度，实现了 sort.Interface 接口
func (s *NamespaceStore) Len() int {
	return len(s.Items)
}

// Less 比较切片中两个元素的大小，实现了 sort.Interface 接口
func (s *NamespaceStore) Less(i, j int) bool {
	return s.Items[i].Priority > s.Items[j].Priority
}

// Swap 交换切片中两个元素的位置，实现了 sort.Interface 接口
func (s *NamespaceStore) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

// -----

func (s *NamespaceStore) Registry(v Object) {
	// fmt.Println("12134", s.Items, s)
	//第一次调用时没有值，所以初始化一个出来，从注册实例中找到这些时
	obj := NewObjectWrapper(v)
	// s.Items = obj
	//这个obj只是返回了值，并没有把obj赋值给s.Items
	//判断name和version一致不
	old, index := s.getWithIndex(obj.Name, obj.Version)
	//1,当切片没有这个的时候，把服务追加进去
	//2.当有重负的时候，返回为空，然后追加进去替换，如果AllowOverwrite不等于true，就会报错
	// 没有, 直接添加
	if old == nil {
		s.Items = append(s.Items, obj)
		return
	}

	// 有, 允许盖写则直接修改
	if obj.AllowOverwrite {
		s.setWithIndex(index, obj)
		return
	}

	// 有, 不允许修改
	panic(fmt.Sprintf("ioc obj %s has registed", obj.Name))
}

// 第一次调用返回为空，第二次才赋值给s.Items
func (s *NamespaceStore) getWithIndex(name, version string) (*ObjectWrapper, int) {
	// fmt.Println("1234", s.Items, s)
	for i := range s.Items {
		// fmt.Println("1234", i, s.Items[i], s.Items)
		obj := s.Items[i] //i=0
		if obj.Name == name && obj.Version == version {
			return obj, i
		}
	}
	return nil, -1
}

// 是否允许同名对象被替换, 默认不允许被替换.根据注册先后
func (s *NamespaceStore) setWithIndex(index int, obj *ObjectWrapper) {
	s.Items[index] = obj
}

func (s *NamespaceStore) List() (uids []string) {
	for i := range s.Items {
		// fmt.Println("i-----", i, s.Items[i])
		item := s.Items[i]
		uids = append(uids, ObjectUid(item))
	}
	return
}

// 注册的名字由来
func ObjectUid(o *ObjectWrapper) string {
	return fmt.Sprintf("%s.%s.%d ", o.Name, o.Version, o.Priority)
}

func (s *NamespaceStore) Get(name string, opts ...GetOption) Object {
	opt := defaultOption().Apply(opts...)
	// opt.version = "v2"
	// fmt.Println("--------", name, opt.version)
	obj, _ := s.getWithIndex(name, opt.version)
	if obj == nil {
		return nil
	}
	return obj.Value
}

// // 对象的优先级
// func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
// 	s.Priority = v
// 	return s
// }

// 根据对象对象加载对象
func (s *NamespaceStore) Load(target any, opts ...GetOption) (Object, error) {
	//根据传进来的any，去反射，判断是什么类型和返回什么值
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)

	var obj Object
	switch t.Kind() {
	//t.Kind() 判断为指针ptr
	case reflect.Interface:
		objs := s.ImplementInterface(t)
		if len(objs) > 0 {
			obj = objs[0]
		}
	default:
		//t.String()=*ioc_test.TestObject
		fmt.Println("t.string:", t.String())
		//load加载，当有name方法时断言会失败
		//ioc-list: [b.v1]
		//t.string: *ioc_test.c
		obj = s.Get(t.String(), opts...)
	}

	// 注入值
	if obj != nil {
		objValue := reflect.ValueOf(obj)
		if !(v.Kind() == reflect.Ptr && objValue.Kind() == reflect.Ptr) {
			return nil, fmt.Errorf("target and object must both be pointers or non-pointers")
		}
		v.Elem().Set(objValue.Elem())
	}
	return obj, nil
}

// 寻找实现了接口的对象
func (s *NamespaceStore) ImplementInterface(objType reflect.Type, opts ...GetOption) (objs []Object) {
	opt := defaultOption().Apply(opts...)

	for i := range s.Items {
		fmt.Println("-----------", i, s.Items[i])
		o := s.Items[i]
		// 断言获取的对象是否满足接口类型
		if o != nil && reflect.TypeOf(o.Value).Implements(objType) {
			if o.Version == opt.version {
				objs = append(objs, o.Value)
			}
		}
	}
	return
}

// 对象个数统计
func (s *NamespaceStore) Count() int {
	return len(s.Items)
}

// func (s *NamespaceStore) Len() int {
// 	return len(s.Items)
// }

// 这个是给加载配置文件用的.api路由也用了
func (s *NamespaceStore) ForEach(fn func(*ObjectWrapper)) {
	for i := range s.Items {
		item := s.Items[i]
		fn(item)
	}
}

// 销毁对象，此接口，暂未调用
func (s *NamespaceStore) Close(ctx context.Context) error {
	errs := []string{}
	for i := range s.Items {
		obj := s.Items[i]
		if err := obj.Value.Close(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("close error, %s", strings.Join(errs, ","))
	}
	return nil
}
