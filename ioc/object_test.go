package ioc_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/datasource"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestObjectLoad(t *testing.T) {
	/*
		将*TestObject对象 注册到默认空间
		通过环境变量配置TestObject对象
		os.Setenv("ATTR1", "a1")
		os.Setenv("ATTR2", "a2")
		暂未开发用处
	*/

	// 除了采用Get直接获取对象, 也可以通过Load动态加载, 等价于获取后赋值
	// m := datasource.DB(context.Background())
	// t.Log(m)
	fmt.Println("ioc-list:", ioc.Config().List())
	fmt.Println("ioc-count:", ioc.Config().Count())
	ioc.Default().Registry(&c{}) //b.v1
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "../etc/application.toml"
	ioc.DevelopmentSetup(req)
	cc := ioc.Default().Get("b").(*c)
	fmt.Println("--------1", cc.Version(), cc.Hello())
	ioc.Api().Registry(&c{})
	r := gin.Default()

	// 加载 Gin API，指定路径前缀和路由器对象
	// ioc.LoadGinApi("/api/v1", r)

	// 启动 Gin 服务器
	r.Run(":8080")
	// obj := &c{}
	// _, err = ioc.Default().Load(obj)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // a := cc.(*c)
	// fmt.Println(obj)
}

// func init() {
// 	os.Setenv("DATASOURCE_HOST", "35.220.136.110")
// 	os.Setenv("DATASOURCE_PORT", "3306")
// 	os.Setenv("DATASOURCE_DB", "sso")
// 	os.Setenv("DATASOURCE_USERNAME", "root")
// 	os.Setenv("DATASOURCE_PASSWORD", "cys000522")
// }

type c struct {
	// Attr1 string `toml:"attr1" env:"ATTR1"`
	// Attr2 string `toml:"attr2" env:"ATTR2"`
	ioc.ObjectImpl
	db *gorm.DB
}

func (t *c) Hello() error {
	t.db = datasource.DB(context.Background())
	// err := t.db.Model(&Users{}).Where("id = 37").Delete(&Users{}).Error
	return nil
}

func (t *c) Hello1() string {
	return "this is test for v1"
}

type TestService interface {
	Hello() string
}

func (i *c) Init() error {
	return nil
}

func (i *c) Name() string {
	return "b"
}

const (
	DEFAULT_VERSIO = "v1"
)

func (i *c) Version() string {
	// return ioc.DEFAULT_VERSION
	return DEFAULT_VERSIO
}

func (i *c) Priority() int {
	return 8
}

func (i *c) AllowOverwrite() bool {
	return false
}

func (i *c) Close(ctx context.Context) error {
	return nil
}

func (m *c) Registry(router gin.IRouter) {
	// 在这里注册 Gin 路由
	router.GET("/endpoint", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Gin API"})
	})

	router.GET("/endpoint1", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from Gin API"})
	})
}
