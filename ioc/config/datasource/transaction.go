package datasource

import (
	"context"

	// GORM。
	"gorm.io/gorm"
)

// GetTransactionOrDB 从 Context 获取事务。
// 如果 Context 中存在事务：返回事务 tx
// 如果 Context 中不存在事务：返回默认 Primary。
// 注意：
// 如果默认 DB 注册了 DBResolver，直接使用 m.db.WithContext(ctx) 后，GORM Resolver 仍然会根据操作类型执行读写分离。
func (m *dataSource) GetTransactionOrDB(ctx context.Context) *gorm.DB {
	// 从 Context 中查找事务。
	tx := GetTransactionFromCtx(ctx)
	// 如果存在事务，则直接返回事务。
	if tx != nil {
		return tx
	}
	// 如果 datasource 尚未初始化，则返回 nil。
	if m.db == nil {
		return nil
	}
	// 没有事务时使用默认数据库。
	return m.db.WithContext(ctx)
}

// GetTransactionOrDBByName 从 Context 获取指定数据库事务。
// 如果 Context 已经存在事务：优先返回事务。
// 如果 Context 没有事务：返回指定数据库。
func (m *dataSource) GetTransactionOrDBByName(ctx context.Context, name string) *gorm.DB {
	// 首先尝试从 Context 获取事务。
	tx := GetTransactionFromCtx(ctx)
	// 已经存在事务时优先使用事务。
	if tx != nil {
		return tx
	}
	// 获取指定独立数据库。
	db := m.DBByName(name)
	// 数据库不存在时返回 nil。
	if db == nil {
		return nil
	}
	// 将 Context 绑定到 GORM DB。
	return db.WithContext(ctx)
}

// GetTransactionFromCtx 从 Context 中获取事务。
// 事务保存时使用 TransactionCtxKey{} 作为 key。
func GetTransactionFromCtx(ctx context.Context) *gorm.DB {
	// Context 为空时不可能存在事务。
	if ctx == nil {
		return nil
	}
	// 从 Context 中读取事务。
	tx, ok := ctx.Value(TransactionCtxKey{}).(*gorm.DB)
	// 类型正确时返回事务。
	if !ok {
		return nil
	}
	return tx
}

// TransactionCtxKey 是事务 Context 使用的 key。
// 使用独立 struct 类型作为 key，
// 可以避免与其他 package 的 Context key 冲突。
type TransactionCtxKey struct{}
