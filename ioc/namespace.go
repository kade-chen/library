package ioc

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/kade-chen/library/tools/file"
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// 调用它返回namespace 和 items
var (
	store = newDefaultStore()
)

type defaultStore struct {
	conf  *LoadConfigRequest
	store []*NamespaceStore
}

// 调用store会默认调用这个方法
// 此处大的优先
func newDefaultStore() *defaultStore {
	return &defaultStore{
		store: []*NamespaceStore{
			newNamespaceStore(CONFIG_NAMESPACE).SetPriority(99),
			newNamespaceStore(DEFAULT_NAMESPACE).SetPriority(9),
			newNamespaceStore(CONTROLLER_NAMESPACE).SetPriority(0),
			newNamespaceStore(API_NAMESPACE).SetPriority(-99),
		},
	}
}

// 获取一个对象存储空间，items和namespace的切片
func (s *defaultStore) Namespace(namespace string) *NamespaceStore {
	//循环判断传进来的nacespace和newDefaultStore()初始化时候的namespace是否一样
	for i := range s.store {
		item := s.store[i]
		// fmt.Println("snsnnsns", item)
		//item 为空，ObjectWrapper，注册前调用时为空
		if item.Namespace == namespace {
			return item
		}
	}
	//如果NamespaceStore不存在这个namespace则创建一个新的namespace
	ns := newNamespaceStore(namespace)
	s.store = append(s.store, ns)
	return ns
}

func newNamespaceStore(namespace string) *NamespaceStore {
	return &NamespaceStore{
		Namespace: namespace,
		Items:     []*ObjectWrapper{},
	}
}

// 对象的优先级
func (s *NamespaceStore) SetPriority(v int) *NamespaceStore {
	s.Priority = v
	return s
}

// ------------------load  InitIocObject------------------
// Here are the methods that load requires
// 初始化托管的所有对象
func (s *defaultStore) InitIocObject() error {
	s.Sort()

	for i := range s.store {
		item := s.store[i]
		//正在的初始化，在这里做的
		err := item.Init()
		if err != nil {
			return fmt.Errorf("[%s] %s", item.Namespace, err)
		}
	}
	return nil
}

// 根据对象的优先级进行排序
// 想要用sort的interface methods就必须实现接口的方法len swap less
func (s *defaultStore) Sort() {
	sort.Sort(s)
}

// Len 返回切片的长度，实现了 sort.Interface 接口
func (s *defaultStore) Len() int {
	return len(s.store)
}

// Less 比较切片中两个元素的大小，实现了 sort.Interface 接口
func (s *defaultStore) Less(i, j int) bool {
	return s.store[i].Priority > s.store[j].Priority
}

// Swap 交换切片中两个元素的位置，实现了 sort.Interface 接口
func (s *defaultStore) Swap(i, j int) {
	s.store[i], s.store[j] = s.store[j], s.store[i]
}

// 注册的实例的初始化，都是在这做的
func (s *NamespaceStore) Init() error {
	s.Sort()
	for i := range s.Items {
		obj := s.Items[i]
		// fmt.Println("------init对象:", obj)
		err := obj.Value.Init()
		if err != nil {
			return fmt.Errorf("init object %s error, %s", obj.Name, err)
		}
	}
	return nil
}

// ------------------load  LoadConfig------------------
// 加载对象配置
func (s *defaultStore) LoadConfig(req *LoadConfigRequest) error {
	errs := []string{}

	// 优先加载环境变量
	if req.ConfigEnv.Enabled {
		for i := range s.store {
			item := s.store[i]
			//使用env库解析环境变量w.Value，并将其存储到相应的结构体字段中，同时只关注那些以指定前缀开头的环境变量。
			err := item.LoadFromEnv(req.ConfigEnv.Prefix)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	// 再加载配置文件
	if req.ConfigFile.Enabled {
		for i := range s.store {
			item := s.store[i]
			err := item.LoadFromFile(req.ConfigFile.Path)
			if err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}

	s.conf = req
	return nil
}

// 从环境变量中加载对象配置
func (i *NamespaceStore) LoadFromEnv(prefix string) error {
	errs := []string{}
	i.ForEach(func(w *ObjectWrapper) {
		err := env.Parse(w.Value, env.Options{
			Prefix: prefix,
		})
		if err != nil {
			errs = append(errs, err.Error())
		}
	})
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}

	return nil
}

// 从环境配置文件中加载对象配置
func (i *NamespaceStore) LoadFromFile(filename string) error {
	if filename == "" {
		return nil
	}

	fileType := filepath.Ext(filename)
	if err := ValidateFileType(fileType); err != nil {
		return err
	}

	// 准备一个map读取配置
	cfg := map[string]any{}
	i.ForEach(func(w *ObjectWrapper) {
		cfg[w.Value.Name()] = w.Value
	})

	var err error
	switch fileType {
	case ".toml":
		_, err = toml.DecodeFile(filename, &cfg)
	case ".yml", ".yaml":
		err = file.ReadYamlFile(filename, &cfg)
	case ".json":
		err = file.ReadJsonFile(filename, &cfg)
	default:
		err = fmt.Errorf("unspport format: %s", fileType)
	}
	if err != nil {
		return err
	}

	// 加载到对象中
	// fmt.Println("1234", cfg)
	errs := []string{}
	i.ForEach(func(w *ObjectWrapper) {
		// fmt.Println(w.Value, 1)
		//给创建的map对象序列化
		dj, err := json.Marshal(cfg[w.Value.Name()])
		if err != nil {
			errs = append(errs, err.Error())
		}
		// fmt.Println(w.Value, 3)
		//把序列化之后的对象塞进w.value里，这样就读到配置文件中的数据了
		err = json.Unmarshal(dj, w.Value)
		// fmt.Println(w.Value, 2)
		if err != nil {
			errs = append(errs, err.Error())
		}
		// &{{} <nil> 127.0.0.1 8080 api  30 60 60 300 16kb true false 0 <nil> <nil> map[gin:0x140004bc2b8 go-restful:0x140004bc2b0] map[gin:0 go-restful:0]} 0 map[api_doc_path:/swagger_docs enable:true enable_api_doc:true host:127.0.0.1 port:8040 web_framework:go-restful]
		//0之前是w.value的值，0之后是cfg[w.Value.Name()的值➕配置文件序列化的值
		//&{{} <nil> 127.0.0.1 8080 api  30 60 60 300 16kb true false 0 <nil> <nil> map[gin:0x14000070728 go-restful:0x14000070720] map[gin:0 go-restful:0]} 1
		// &{{} <nil> 127.0.0.1 8080 api  30 60 60 300 16kb true false 0 <nil> <nil> map[gin:0x14000070728 go-restful:0x14000070720] map[gin:0 go-restful:0]} 3
		// &{{} 0x1400048c0bd 127.0.0.1 8040 api go-restful 30 60 60 300 16kb true false 0 <nil> <nil> map[gin:0x14000070728 go-restful:0x14000070720] map[gin:0 go-restful:0]} 2
		//最后一步才塞进去的
	})
	if len(errs) > 0 {
		return fmt.Errorf("load config error, %s", strings.Join(errs, ","))
	}
	return nil
}

func ValidateFileType(ext string) error {
	exist := false
	validateFileType := []string{".toml", ".yml", ".yaml", ".json"}
	for _, ft := range validateFileType {
		if ext == ft {
			exist = true
		}
	}
	if !exist {
		return fmt.Errorf("not support format: %s", ext)
	}

	return nil
}

// ------------------Autowire 依赖自动注入-----------------
func (s *defaultStore) Autowire() error {
	for i := range s.store {
		item := s.store[i]
		err := item.Autowire()
		if err != nil {
			return fmt.Errorf("[%s] %s", item.Namespace, err)
		}
	}
	return nil
}

func (i *NamespaceStore) Autowire() error {
	i.ForEach(func(w *ObjectWrapper) {
		// 获取对象的类型信息
		pt := reflect.TypeOf(w.Value).Elem()
		// 获取对象的值信息
		v := reflect.ValueOf(w.Value).Elem()

		// 遍历对象的字段
		for i := 0; i < pt.NumField(); i++ {
			// 获取字段的ioc标签
			tag := ParseInjectTag(pt.Field(i).Tag.Get("ioc"))
			if tag.Autowire {
				// 获取字段的类型
				fieldType := v.Field(i).Type()
				var obj Object

				// 根据字段的类型进行不同的处理
				switch fieldType.Kind() {
				case reflect.Interface:
					// 如果字段是接口类型，从命名空间获取实现该接口的对象
					objs := store.Namespace(tag.Namespace).ImplementInterface(fieldType)
					if len(objs) > 0 {
						obj = objs[0]
					}
				default:
					// 如果字段是结构体变量类型，从命名空间获取具体的对象
					if tag.Name == "" {
						tag.Name = fieldType.String()
					}
					obj = store.Namespace(tag.Namespace).Get(tag.Name, WithVersion(tag.Version))
				}

				// 如果成功获取到对象，将其注入到字段中
				if obj != nil {
					v.Field(i).Set(reflect.ValueOf(obj))
				}
			}
		}
	})
	return nil
}

func WithVersion(v string) GetOption {
	return func(o *option) {
		o.version = v
	}
}
