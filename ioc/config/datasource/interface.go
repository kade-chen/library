package datasource

import (
	"context"

	"github.com/kade-chen/library/ioc"
	"gorm.io/gorm"
)

const (
	AppName = "datasource"
)

// transaction is nil
// 此处已经断言成db了，可以直接使用
// func DB(context.Context) *gorm.DB {
// 	// return ioc.Config().Get(AppName).(*dataSource).db.WithContext(context.Background())  这三个功能一样，都是断言后，1.带有上下文 2.带有db的数据库连接，1，3一样
// 	// return ioc.Config().Get(AppName).(*dataSource).db
// 	return ioc.Config().Get(AppName).(*dataSource).GetTransactionOrDB(context.Background())
// }


func DB() *gorm.DB {
	return Get().db
}

// 从上下文中获取事物, 如果获取不到则直接返回 无事物的DB对象
func DBFromCtx(ctx context.Context) *gorm.DB {
	return Get().GetTransactionOrDB(ctx)
}

func Get() *dataSource {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*dataSource)
}
