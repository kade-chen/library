package datasource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/datasource"
)

func TestMysql(t *testing.T) {
	fmt.Println("ioc-list:", ioc.Config().List())
	fmt.Println("ioc-count:", ioc.Config().Count())
	m := datasource.DB(context.Background())
	t.Log(m)
}
