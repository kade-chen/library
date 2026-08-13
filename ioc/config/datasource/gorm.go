package datasource

import (
	"context"
	"fmt"

	// SQLite 驱动。
	"github.com/glebarez/sqlite"
	// Vault Lease revoke 所需要的请求结构。
	"github.com/hashicorp/vault-client-go/schema"
	// IOC 注册。
	"github.com/kade-chen/library/ioc"
	// 获取应用名称。
	"github.com/kade-chen/library/ioc/config/application"
	// 日志。
	"github.com/kade-chen/library/ioc/config/log"
	// Trace 配置。
	"github.com/kade-chen/library/ioc/config/trace"
	// Vault Client。
	"github.com/kade-chen/library/ioc/config/vault"
	// Zerolog Logger。
	"github.com/rs/zerolog"
	// GORM OpenTelemetry 插件。
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	// MySQL 驱动。
	"gorm.io/driver/mysql"
	// PostgreSQL 驱动。
	"gorm.io/driver/postgres"
	// GORM。
	"gorm.io/gorm"
	// GORM 多数据库、读写分离插件。
	"gorm.io/plugin/dbresolver"
)

// init 在 datasource package 被加载时执行。
// 将默认 datasource 注册到 IOC。
func init() {
	ioc.Config().Registry(defaultConfig)
}

// defaultConfig 是 datasource 的默认配置。
// 如果外部没有配置 datasource，则使用这里的默认值。
var defaultConfig = &dataSource{
	// 默认数据库类型为 MySQL。
	Provider: PROVIDER_MYSQL,
	// 默认数据库地址。
	Host: "127.0.0.1",
	// MySQL 默认端口。
	Port: 3306,
	// 默认使用应用名称作为数据库名称。
	DB: application.Get().Name(),
	// 默认不开启 AutoMigrate。
	AutoMigrate: false,
	// 默认关闭 GORM Debug。
	Debug: false,
	// 默认开启 Trace。
	Trace: true,
	// 默认使用静态账号密码。
	CredentialMode: CREDENTIAL_MODE_STATIC,
	// Vault 用户名字段默认名称。
	VaultUsernameField: "username",
	// Vault 密码字段默认名称。
	VaultPasswordField: "password",
	// 默认自动续期 Vault 动态凭证。
	VaultAutoRenew: true,
	// Vault 租约达到 80% 时开始续期。
	VaultRenewThreshold: 0.8,
	// GORM 默认事务行为。
	SkipDefaultTransaction: false,
	// 默认不开启 DryRun。
	DryRun: false,
	// 默认开启 Prepared Statement。
	PrepareStmt: true,
}

// DatabaseConfig 表示一个数据库配置。
// 一个 DatabaseConfig 表示一个 Primary。
// Primary 下面可以配置多个 Replicas。
type DatabaseConfig struct {
	// 数据库类型。
	//
	// 支持：
	// postgres
	// mysql
	// sqlite
	Provider PROVIDER `json:"provider" yaml:"provider" toml:"provider"`
	// 数据库地址。
	Host string `json:"host" yaml:"host" toml:"host"`
	// 数据库端口。
	Port int `json:"port" yaml:"port" toml:"port"`
	// 数据库名称。
	DB string `json:"database" yaml:"database" toml:"database"`
	// 数据库用户名。
	Username string `json:"username" yaml:"username" toml:"username"`
	// 数据库密码。
	Password string `json:"password" yaml:"password" toml:"password"`
	// 是否开启 Debug。
	Debug bool `json:"debug" yaml:"debug" toml:"debug"`
	// 是否开启 Trace。
	Trace bool `json:"trace" yaml:"trace" toml:"trace"`
	// 当前 Primary 对应的 Replica 列表。
	//
	// 例如：
	//
	// [datasource.databases.cc.replicas.slave1]
	//
	// [datasource.databases.cc.replicas.slave2]
	Replicas map[string]DatabaseConfig `json:"replicas" yaml:"replicas" toml:"replicas"`
}

// dataSource 是整个 datasource 的 IOC 对象。
type dataSource struct {
	// IOC ObjectImpl。
	ioc.ObjectImpl
	// 默认数据库类型。
	Provider PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"PROVIDER"`
	// 默认 Primary 地址。
	Host string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	// 默认 Primary 端口。
	Port int `json:"port" yaml:"port" toml:"port" env:"PORT"`
	// 默认 Primary 数据库名称。
	DB string `json:"database" yaml:"database" toml:"database" env:"DB"`
	// 默认 Primary 用户名。
	Username string `json:"username" yaml:"username" toml:"username" env:"USERNAME"`
	// 默认 Primary 密码。
	Password string `json:"password" yaml:"password" toml:"password" env:"PASSWORD"`
	// 是否执行 AutoMigrate。
	AutoMigrate bool `json:"auto_migrate" yaml:"auto_migrate" toml:"auto_migrate" env:"AUTO_MIGRATE"`
	// 是否开启 GORM Debug。
	Debug bool `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG1"`
	// 是否开启 OpenTelemetry Trace。
	Trace bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`
	// 数据库凭证获取模式。
	CredentialMode CREDENTIAL_MODE `json:"credential_mode" yaml:"credential_mode" toml:"credential_mode" env:"CREDENTIAL_MODE"`
	// Vault Secret 路径。
	VaultPath string `json:"vault_path" yaml:"vault_path" toml:"vault_path" env:"VAULT_PATH"`
	// Vault 用户名字段名称。
	VaultUsernameField string `json:"vault_username_field" yaml:"vault_username_field" toml:"vault_username_field" env:"VAULT_USERNAME_FIELD"`
	// Vault 密码字段名称。
	VaultPasswordField string `json:"vault_password_field" yaml:"vault_password_field" toml:"vault_password_field" env:"VAULT_PASSWORD_FIELD"`
	// 是否自动续期 Vault 动态凭证。
	VaultAutoRenew bool `json:"vault_auto_renew" yaml:"vault_auto_renew" toml:"vault_auto_renew" env:"VAULT_AUTO_RENEW"`
	// Vault 租约续期阈值。
	VaultRenewThreshold float64 `json:"vault_renew_threshold" yaml:"vault_renew_threshold" toml:"vault_renew_threshold" env:"VAULT_RENEW_THRESHOLD"`
	// GORM 是否跳过默认事务。
	SkipDefaultTransaction bool `toml:"skip_default_transaction" json:"skip_default_transaction" yaml:"skip_default_transaction" env:"SKIP_DEFALT_TRANSACTION"`
	// 是否保存所有关联对象。
	FullSaveAssociations bool `toml:"full_save_associations" json:"full_save_associations" yaml:"full_save_associations" env:"FULL_SAVE_ASSOCIATIONS"`
	// 是否只生成 SQL，不真正执行。
	DryRun bool `toml:"dry_run" json:"dry_run" yaml:"dry_run" env:"DRY_RUN"`
	// 是否开启 Prepared Statement。
	PrepareStmt bool `toml:"prepare_stmt" json:"prepare_stmt" yaml:"prepare_stmt" env:"PREPARE_STMT"`
	// 是否关闭 GORM 自动 Ping。
	DisableAutomaticPing bool `toml:"disable_automatic_ping" json:"disable_automatic_ping" yaml:"disable_automatic_ping" env:"DISABLE_AUTOMATIC_PING"`
	// Migration 时是否关闭外键约束创建。
	DisableForeignKeyConstraintWhenMigrating bool `toml:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating" yaml:"disable_foreign_key_constraint_when_migrating" env:"DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING"`
	// Migration 时是否忽略关系。
	IgnoreRelationshipsWhenMigrating bool `toml:"ignore_relationships_when_migrating" json:"ignore_relationships_when_migrating" yaml:"ignore_relationships_when_migrating" env:"IGNORE_RELATIONSHIP_WHEN_MIGRATING"`
	// 是否禁止嵌套事务。
	DisableNestedTransaction bool `toml:"disable_nested_transaction" json:"disable_nested_transaction" yaml:"disable_nested_transaction" env:"DISABLE_NESTED_TRANSACTION"`
	// 是否允许全表 Update。
	AllowGlobalUpdate bool `toml:"allow_global_update" json:"allow_global_update" yaml:"allow_global_update" env:"ALL_GLOBAL_UPDATE"`
	// 是否查询所有字段。
	QueryFields bool `toml:"query_fields" json:"query_fields" yaml:"query_fields" env:"QUERY_FIELDS"`
	// Create 批量插入大小。
	CreateBatchSize int `toml:"create_batch_size" json:"create_batch_size" yaml:"create_batch_size" env:"CREATE_BATCH_SIZE"`
	// 是否开启 GORM 错误转换。
	TranslateError bool `toml:"translate_error" json:"translate_error" yaml:"translate_error" env:"TRANSLATE_ERROR"`
	// 默认 Primary 的 Replica。
	// 例如：[datasource.replicas.bb]
	// [datasource.replicas.bb2]
	Replicas map[string]DatabaseConfig `json:"replicas" yaml:"replicas" toml:"replicas"`
	// 独立数据库。
	// 每个独立数据库自己拥有 Primary + Replicas。
	// 例如：[datasource.databases.cc]
	// [datasource.databases.cc.replicas.slave1]
	Databases map[string]DatabaseConfig `json:"databases" yaml:"databases" toml:"databases"`
	// 默认 Primary 的 GORM DB。
	db *gorm.DB
	// 独立数据库的 Primary DB。
	dbs map[string]*gorm.DB
	// datasource 日志对象。
	log *zerolog.Logger
	// Vault 动态凭证 Lease ID。
	leaseID string
	// Vault Lease 生命周期。
	leaseDuration int64
	// Vault 自动续期停止信号。
	stopRenewal chan struct{}
}

// Name 返回 IOC Bean 名称。
func (m *dataSource) Name() string {
	return AppName
}

// Priority 返回 IOC 初始化优先级。
func (i *dataSource) Priority() int {
	return 999
}

// Init 初始化 datasource。
func (m *dataSource) Init() error {
	// 创建 datasource 专属 logger。
	m.log = log.Sub(m.Name())
	// 设置默认配置。
	if err := m.setDefaults(); err != nil {
		return err
	}
	// 根据 CredentialMode 获取数据库账号密码。
	if err := m.loadCredentials(); err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}
	// ============================================================
	// 1. 初始化默认 Primary
	// ============================================================
	// 使用默认数据库 Dialector 创建 GORM DB。
	db, err := gorm.Open(
		m.Dialector(),
		m.gormConfig(),
	)
	if err != nil {
		return err
	}
	// 如果开启 Trace，则给 GORM 安装 OpenTelemetry 插件。
	if err := m.installTrace(db); err != nil {
		return err
	}
	// 如果开启 Debug，则使用 GORM Debug 模式。
	if m.Debug {
		db = db.Debug()
	}
	// 保存默认 Primary DB。
	m.db = db
	// ============================================================
	// 2. 注册默认 Primary 的 Replicas
	// ============================================================
	// 如果配置了 Replica，则注册 DBResolver。
	if len(m.Replicas) > 0 {
		// 创建 Replica Dialector 列表。
		replicas := make(
			[]gorm.Dialector,
			0,
			len(m.Replicas),
		)
		// 遍历默认数据库的所有 Replica。
		for name, config := range m.Replicas {
			// Replica 名称不能为空。
			if name == "" {
				return fmt.Errorf("datasource replica name cannot be empty")
			}
			// primary 是保留名称。
			if name == "primary" {
				return fmt.Errorf("datasource replica name %q is reserved", name)
			}
			// 没有配置的字段继承默认 Primary 配置。
			config = m.mergeDatabaseConfig(config)
			// 根据 Replica 配置创建 Dialector。
			dialector, err := m.databaseDialector(config)
			if err != nil {
				return fmt.Errorf("create replica %q: %w", name, err)
			}
			// 将 Replica Dialector 加入 Resolver。
			replicas = append(replicas, dialector)
			// 记录 Replica 注册日志。
			m.log.Info().
				Str("name", name).
				Str("provider", string(config.Provider)).
				Str("host", config.Host).
				Int("port", config.Port).
				Str("database", config.DB).
				Msg("datasource replica registered")
		}
		// 注册 DBResolver。
		// Query 默认使用 Replica。
		// Create / Update / Delete 默认使用 Primary。
		// RandomPolicy 会在多个 Replica 之间随机选择。
		err = db.Use(dbresolver.Register(dbresolver.Config{Replicas: replicas, Policy: dbresolver.RandomPolicy{}}).SetMaxIdleConns(100).SetMaxOpenConns(200))
		if err != nil {
			return fmt.Errorf("register primary db resolver: %w", err)
		}
	}
	// 记录默认 Primary 和 Replica 数量。
	m.log.Info().Msgf("datasource primary registered, replicas=%d", len(m.Replicas))
	// ============================================================
	// 3. 初始化独立数据库
	// ============================================================
	// 创建独立数据库 map。
	m.dbs = make(map[string]*gorm.DB, len(m.Databases))
	// 遍历所有独立数据库。
	for name, config := range m.Databases {
		// 数据库名称不能为空。
		if name == "" {
			return fmt.Errorf("datasource database name cannot be empty")
		}
		// primary 是保留名称。
		if name == "primary" {
			return fmt.Errorf("datasource database name %q is reserved", name)
		}
		// 初始化独立数据库。
		database, err := m.openIndependentDatabase(name, config)
		if err != nil {
			return fmt.Errorf("open datasource database %q: %w", name, err)
		}
		// 保存独立数据库 Primary。
		m.dbs[name] = database
		// 记录独立数据库注册日志。
		m.log.Info().
			Str("name", name).
			Str("provider", string(config.Provider)).
			Str("host", config.Host).
			Int("port", config.Port).
			Str("database", config.DB).
			Msgf(
				"datasource database primary registered, replicas=%d",
				len(config.Replicas),
			)
	}
	// ============================================================
	// 4. 启动 Vault 动态凭证自动续期
	// ============================================================
	// 只有 Vault Dynamic 模式并且开启 AutoRenew 时才启动续期。
	if m.CredentialMode.NeedsRenewal() && m.VaultAutoRenew {
		go m.startCredentialRenewal()
	}
	return nil
}

// openIndependentDatabase 初始化一个独立数据库。
// 一个独立数据库同样拥有：
//
//	Primary
//	    ↓
//	Replicas 因此 cc 和默认 datasource 可以拥有完全独立的主从结构。
func (m *dataSource) openIndependentDatabase(name string, config DatabaseConfig) (*gorm.DB, error) {
	// 将没有填写的配置继承默认 datasource 配置。
	config = m.mergeDatabaseConfig(config)
	// 创建独立数据库 Primary 的 Dialector。
	dialector, err := m.databaseDialector(config)
	if err != nil {
		return nil, err
	}
	// 打开独立数据库 Primary。
	db, err := gorm.Open(
		dialector,
		m.gormConfig(),
	)
	if err != nil {
		return nil, err
	}
	// 给独立数据库安装 Trace。
	if err := m.installTrace(db); err != nil {
		return nil, err
	}
	// 独立数据库开启 Debug 时使用 Debug 模式。
	if m.Debug || config.Debug {
		db = db.Debug()
	}
	// ============================================================
	// 独立数据库自己的 Replicas
	// ============================================================
	// 如果独立数据库配置了 Replica，则注册自己的 Resolver。
	if len(config.Replicas) > 0 {
		// 创建独立数据库 Replica Dialector 列表。
		replicas := make([]gorm.Dialector, 0, len(config.Replicas))
		// 遍历独立数据库 Replica。
		for replicaName, replicaConfig := range config.Replicas {
			// 继承独立数据库 Primary 的默认配置。
			// 例如 slave1 只写 host：host = "10.0.1.3"
			// database、username、password 可以从 cc Primary 继承。
			replicaConfig = mergeDatabaseConfigFrom(config, replicaConfig)
			// 创建 Replica Dialector。
			dialector, err := m.databaseDialector(replicaConfig)
			if err != nil {
				return nil, fmt.Errorf("create database %s replica %s: %w", name, replicaName, err)
			}
			// 加入 Replica 列表。
			replicas = append(replicas, dialector)
			// 记录独立数据库 Replica 注册日志。
			m.log.Info().
				Str("database", name).
				Str("replica", replicaName).
				Str("provider", string(replicaConfig.Provider)).
				Str("host", replicaConfig.Host).
				Int("port", replicaConfig.Port).
				Str("database_name", replicaConfig.DB).
				Msg("datasource database replica registered")
		}
		// 给独立数据库注册 DBResolver。
		// cc 查询：SELECT -> cc replicas
		// cc 写入： INSERT / UPDATE / DELETE -> cc Primary
		err = db.Use(dbresolver.Register(dbresolver.Config{Replicas: replicas, Policy: dbresolver.RandomPolicy{}}).SetMaxIdleConns(100).SetMaxOpenConns(200))
		if err != nil {
			return nil, fmt.Errorf("register database %s resolver: %w", name, err)
		}
	}
	return db, nil
}

// mergeDatabaseConfig 将 Replica 配置中没有填写的字段继承默认 datasource。
func (m *dataSource) mergeDatabaseConfig(config DatabaseConfig) DatabaseConfig {
	// 没配置 Provider，则继承默认 Provider。
	if config.Provider == "" {
		config.Provider = m.Provider
	}
	// 没配置 Host，则继承默认 Host。
	if config.Host == "" {
		config.Host = m.Host
	}
	// 没配置 Port，则继承默认 Port。
	if config.Port == 0 {
		config.Port = m.Port
	}
	// 没配置 Username，则继承默认 Username。
	if config.Username == "" {
		config.Username = m.Username
	}
	// 没配置 Password，则继承默认 Password。
	if config.Password == "" {
		config.Password = m.Password
	}
	// 没配置数据库名称，则继承默认数据库名称。
	if config.DB == "" {
		config.DB = m.DB
	}
	return config
}

// mergeDatabaseConfigFrom 将 Replica 配置中没有填写的字段
// 继承它所属的 Primary 配置。
// 例如：
//
//	cc Primary:
//	    host = 10.0.1.2
//	    port = 5432
//	    database = cc
//	cc slave1:
//	    host = 10.0.1.3
//
// 最终 slave1：
//
//	host      = 10.0.1.3
//	port      = 5432
//	database  = cc
func mergeDatabaseConfigFrom(primary DatabaseConfig, replica DatabaseConfig) DatabaseConfig {
	// Replica 没有 Provider 时继承 Primary。
	if replica.Provider == "" {
		replica.Provider = primary.Provider
	}
	// Replica 没有 Host 时继承 Primary。
	if replica.Host == "" {
		replica.Host = primary.Host
	}
	// Replica 没有 Port 时继承 Primary。
	if replica.Port == 0 {
		replica.Port = primary.Port
	}
	// Replica 没有 DB 时继承 Primary。
	if replica.DB == "" {
		replica.DB = primary.DB
	}
	// Replica 没有 Username 时继承 Primary。
	if replica.Username == "" {
		replica.Username = primary.Username
	}
	// Replica 没有 Password 时继承 Primary。
	if replica.Password == "" {
		replica.Password = primary.Password
	}
	return replica
}

// gormConfig 创建统一的 GORM 配置。
//
// 所有 Primary 和 Replica 都使用相同的 GORM 基础配置。
func (m *dataSource) gormConfig() *gorm.Config {
	return &gorm.Config{
		// 是否关闭 GORM 默认事务。
		SkipDefaultTransaction: m.SkipDefaultTransaction,
		// 是否保存所有关联对象。
		FullSaveAssociations: m.FullSaveAssociations,
		// 是否只生成 SQL。
		DryRun: m.DryRun,
		// 是否开启 Prepared Statement。
		PrepareStmt: m.PrepareStmt,
		// 是否关闭自动 Ping。
		DisableAutomaticPing: m.DisableAutomaticPing,
		// Migration 是否关闭外键约束。
		DisableForeignKeyConstraintWhenMigrating: m.DisableForeignKeyConstraintWhenMigrating,
		// Migration 是否忽略关系。
		IgnoreRelationshipsWhenMigrating: m.IgnoreRelationshipsWhenMigrating,
		// 是否禁止嵌套事务。
		DisableNestedTransaction: m.DisableNestedTransaction,
		// 是否允许全局 Update。
		AllowGlobalUpdate: m.AllowGlobalUpdate,
		// 是否查询全部字段。
		QueryFields: m.QueryFields,
		// Create Batch Size。
		CreateBatchSize: m.CreateBatchSize,
		// 是否开启错误转换。
		TranslateError: m.TranslateError,
		// 使用当前 datasource 的自定义 Logger。
		Logger: newGormCustomLogger(m.log),
	}
}

// installTrace 给 GORM DB 安装 OpenTelemetry Trace 插件。
func (m *dataSource) installTrace(db *gorm.DB) error {
	// 全局 Trace 开启并且 datasource Trace 开启时才安装。
	if trace.Get().Enable && m.Trace {
		// 安装 GORM OpenTelemetry 插件。
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			return err
		}
	}
	return nil
}

// databaseDialector 根据 Provider 创建对应数据库 Dialector。
func (m *dataSource) databaseDialector(config DatabaseConfig) (gorm.Dialector, error) {
	switch config.Provider {
	case PROVIDER_POSTGRES:
		// PostgreSQL DSN。
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			config.Host,
			config.Username,
			config.Password,
			config.DB,
			config.Port,
		)
		return postgres.Open(dsn), nil
	case PROVIDER_SQLITE:
		// SQLite 使用数据库文件路径。
		return sqlite.Open(config.DB), nil
	case PROVIDER_MYSQL:
		// MySQL DSN。
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.DB,
		)
		return mysql.Open(dsn), nil
	default:
		// 不支持的 Provider 直接返回错误。
		return nil, fmt.Errorf("unsupported provider: %q", config.Provider)
	}
}

// Dialector 返回默认 Primary 的 Dialector。
// 保持原有 API，避免影响已有业务代码。
func (m *dataSource) Dialector() gorm.Dialector {
	// 将默认 datasource 转换成 DatabaseConfig。
	config := DatabaseConfig{
		Provider: m.Provider,
		Host:     m.Host,
		Port:     m.Port,
		DB:       m.DB,
		Username: m.Username,
		Password: m.Password,
	}
	// 创建默认 Primary Dialector。
	dialector, _ := m.databaseDialector(config)
	return dialector
}

// DBByName 根据数据库名称获取独立数据库 Primary。
// name 为空：返回默认 Primary。
// name = primary：返回默认 Primary。
// name = master：兼容旧代码，返回默认 Primary。
// 其他名称：返回对应独立数据库。
func (m *dataSource) DBByName(name string) *gorm.DB {
	// 默认数据库。
	if name == "" || name == "primary" || name == "master" {
		return m.db
	}
	// 返回独立数据库 Primary。
	// 注意：这个 DB 已经安装了自己的 DBResolver。
	// 所以： DBByName("cc").First(...) 会自动走 cc Replica。
	return m.dbs[name]
}

// Close 关闭 datasource 所有数据库连接。
func (m *dataSource) Close(ctx context.Context) error {
	// 停止 Vault 自动续期。
	if m.stopRenewal != nil {
		close(m.stopRenewal)
		m.stopRenewal = nil
	}
	// 如果使用 Vault Dynamic Credential，
	// 则主动撤销当前 Lease。
	if m.CredentialMode ==
		CREDENTIAL_MODE_VAULT_DYNAMIC &&
		m.leaseID != "" {
		// 获取 Vault Client。
		vaultClient := vault.Client()
		// Vault Client 存在时执行 Lease revoke。
		if vaultClient != nil {
			if _, err := vaultClient.System.
				LeasesRevokeLease(ctx, schema.LeasesRevokeLeaseRequest{LeaseId: m.leaseID}); err != nil {
				// 记录撤销失败。
				m.log.Error().Err(err).Msgf("failed to revoke lease %s", m.leaseID)
			}
		}
	}
	// 用于保存第一个关闭错误。
	var firstErr error
	// ============================================================
	// 关闭默认 Primary
	// ============================================================
	if m.db != nil {
		// 获取底层 *sql.DB。
		sqlDB, err := m.db.DB()
		if err != nil {
			// 保存错误。
			firstErr = err
		} else if err := sqlDB.Close(); err != nil {
			// 保存关闭错误。
			firstErr = err
		}
	}
	// ============================================================
	// 关闭所有独立数据库 Primary
	// ============================================================
	for name, db := range m.dbs {
		// DB 不存在时跳过。
		if db == nil {
			continue
		}
		// 获取底层 *sql.DB。
		sqlDB, err := db.DB()
		if err != nil {
			// 只保存第一个错误。
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		// 关闭数据库连接池。
		if err := sqlDB.Close(); err != nil {
			// 只保存第一个错误。
			if firstErr == nil {
				firstErr = err
			}
		}
		// 记录关闭日志。
		m.log.Info().Str("name", name).Msg("datasource database closed")
	}
	// 清空独立数据库 map。
	m.dbs = nil
	// 清空默认 Primary。
	m.db = nil
	// 返回关闭过程中出现的第一个错误。
	return firstErr
}

// Version 返回 datasource 版本。
func (i *dataSource) Version() string {
	return ""
}

// AllowOverwrite 表示是否允许 IOC 覆盖当前对象。
func (i *dataSource) AllowOverwrite() bool {
	return false
}
