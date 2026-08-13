package datasource_test

import (
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/datasource"
)

// postgres
func TestMysql(t *testing.T) {
	fmt.Println("111")
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "test/postgres.toml"
	ioc.DevelopmentSetup(req)

	fmt.Println("ioc-list:", ioc.Config().List())
	fmt.Println("ioc-count:", ioc.Config().Count())
	m := datasource.DB()
	t.Log(m)
}
