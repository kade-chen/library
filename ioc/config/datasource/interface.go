package datasource

import (
	"context"
	"fmt"

	// IOC。
	"github.com/kade-chen/library/ioc"

	// GORM。
	"gorm.io/gorm"

	// GORM DBResolver。
	"gorm.io/plugin/dbresolver"
)

// AppName 是 datasource 在 IOC 中注册时使用的名称。
const (
	AppName = "datasource"
)

// DB 获取默认数据库 Primary。
//
// 注意：
// 这里返回的是已经注册 DBResolver 的 *gorm.DB。
//
// 所以：
//
//	SELECT  -> Replica
//	INSERT  -> Primary
//	UPDATE  -> Primary
//	DELETE  -> Primary
//
// GORM DBResolver 官方支持这种自动读写分离。 [oai_citation:1‡GORM](https://gorm.io/docs/dbresolver.html?utm_source=chatgpt.com)
func DB() *gorm.DB {

	// 返回默认 Primary。
	return Get().db
}

// DBByName 获取指定的独立数据库 Primary。
//
// 例如：
//
//	datasource.DBByName("cc")
//
// cc 本身如果配置了 replicas：
//
//	SELECT  -> cc replicas
//	INSERT  -> cc primary
//	UPDATE  -> cc primary
//	DELETE  -> cc primary
func DBByName(
	name string,
) *gorm.DB {

	// 根据名称获取数据库。
	return Get().DBByName(name)
}

// DBRead 强制默认数据库走 Replica。
// 即使当前操作是 First / Find，也明确要求使用 Replica。
func DBRead() *gorm.DB {
	// Read 会强制 DBResolver 使用 Replica。
	return Get().db.Clauses(
		dbresolver.Read,
	)
}

// DBWrite 强制默认数据库走 Primary。
// 即使当前执行的是 SELECT，也会使用 Primary。
// 例如：
//
//	datasource.DBWrite().First(&user)
func DBWrite() *gorm.DB {
	// Write 会强制 DBResolver 使用 Primary。
	return Get().db.Clauses(
		dbresolver.Write,
	)
}

// DBReadByName 强制指定独立数据库走 Replica。
// 例如：datasource.DBReadByName("cc").First(&user)
func DBReadByName(name string) *gorm.DB {
	// 获取指定独立数据库。
	db := Get().DBByName(name)
	// 数据库不存在时返回 nil。
	if db == nil {
		return nil
	}
	// 强制使用 Replica。
	return db.Clauses(
		dbresolver.Read,
	)
}

// DBWriteByName 强制指定独立数据库走 Primary。
// 例如：datasource.DBWriteByName("cc").First(&user)
func DBWriteByName(name string) *gorm.DB {
	// 获取指定独立数据库。
	db := Get().DBByName(name)
	// 数据库不存在时返回 nil。
	if db == nil {
		return nil
	}
	// 强制使用 Primary。
	return db.Clauses(
		dbresolver.Write,
	)
}

// DBFromCtx 从 Context 获取事务。
// 如果 Context 中存在事务，则返回事务。
// 如果不存在，则返回默认数据库。
func DBFromCtx(ctx context.Context) *gorm.DB {
	// 调用 datasource 内部的事务获取方法。
	return Get().GetTransactionOrDB(ctx)
}

// DBFromCtxByName 从 Context 获取指定数据库事务。
// 如果 Context 中存在事务，则优先返回事务。
// 如果没有事务，则返回指定独立数据库。
func DBFromCtxByName(ctx context.Context, name string) *gorm.DB {
	// 调用 datasource 内部方法。
	return Get().GetTransactionOrDBByName(ctx, name)
}

// Get 获取 datasource IOC 实例。
func Get() *dataSource {
	// 从 IOC 中获取 datasource。
	obj := ioc.Config().Get(AppName)
	// 如果 IOC 中不存在，则返回默认配置。
	if obj == nil {
		return defaultConfig
	}
	// 转换成 dataSource。
	return obj.(*dataSource)
}

func (m *dataSource) GetRegisteredDatabases() {
	primaryCount := 0
	replicaCount := 0
	fmt.Println("========== Registered Databases ==========")
	if m.db != nil {
		primaryCount++
		fmt.Println("default")
		fmt.Println("├── primary")
		for name := range m.Replicas {
			replicaCount++
			fmt.Printf("├── replica: %s\n", name)
		}
	}

	for name, db := range m.dbs {
		if db == nil {
			continue
		}

		primaryCount++
		fmt.Printf("%s\n", name)
		fmt.Println("├── primary")

		config := m.Databases[name]
		for replicaName := range config.Replicas {
			replicaCount++
			fmt.Printf("├── replica: %s\n", replicaName)
		}
	}

	fmt.Println("==========================================")
	fmt.Printf(
		"Summary: primary=%d, replica=%d, total=%d\n",
		primaryCount,
		replicaCount,
		primaryCount+replicaCount,
	)
}
