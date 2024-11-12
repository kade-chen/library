package datasource

import (
	"context"

	"gorm.io/gorm"
)

// 从上下文中获取事物, 如果获取不到则直接返回 无事物的DB对象
func (m *dataSource) GetTransactionOrDB(ctx context.Context) *gorm.DB {
	// 从上下文中获取数据库事务
	db := GetTransactionFromCtx(ctx)
	// 如果事务存在，则返回该事务
	if db != nil {
		return db
	}
	// 否则，返回默认的数据库连接（非事务）
	return m.db.WithContext(ctx)
}

// 从上下文中获取数据库事务
func GetTransactionFromCtx(ctx context.Context) *gorm.DB {
	if ctx != nil {
		// 通过上下文键 TransactionCtxKey 获取事务对象
		tx, ok := ctx.Value(TransactionCtxKey{}).(*gorm.DB)
		if ok {
			return tx
		}
	}
	// 如果上下文中没有事务对象，则返回 nil
	return nil
}

// 上下文键的类型，用于在上下文中存储事务对象
type TransactionCtxKey struct{}
