package datasource

import (
	"context"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/hashicorp/vault-client-go/schema"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/application"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/kade-chen/library/ioc/config/trace"
	"github.com/kade-chen/library/ioc/config/vault"
	"github.com/rs/zerolog"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &dataSource{
	Provider:    PROVIDER_MYSQL,
	Host:        "127.0.0.1",
	Port:        3306,
	DB:          application.Get().Name(),
	AutoMigrate: false,
	Debug:       false,
	Trace:       true,

	// Vault 凭证默认配置
	CredentialMode:      CREDENTIAL_MODE_STATIC,
	VaultUsernameField:  "username",
	VaultPasswordField:  "password",
	VaultAutoRenew:      true,
	VaultRenewThreshold: 0.8,

	SkipDefaultTransaction: false,
	DryRun:                 false,
	PrepareStmt:            true,
}

type dataSource struct {
	ioc.ObjectImpl
	Provider    PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"PROVIDER"`
	Host        string   `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port        int      `json:"port" yaml:"port" toml:"port" env:"PORT"`
	DB          string   `json:"database" yaml:"database" toml:"database" env:"DB"`
	Username    string   `json:"username" yaml:"username" toml:"username" env:"USERNAME"`
	Password    string   `json:"password" yaml:"password" toml:"password" env:"PASSWORD"`
	AutoMigrate bool     `json:"auto_migrate" yaml:"auto_migrate" toml:"auto_migrate" env:"AUTO_MIGRATE"`
	Debug       bool     `json:"debugs" yaml:"debugs" toml:"debugs" env:"DEBUGS"`
	Trace       bool     `toml:"trace" json:"trace" yaml:"trace"  env:"TRACE"`

	// Vault 凭证配置
	CredentialMode CREDENTIAL_MODE `json:"credential_mode" yaml:"credential_mode" toml:"credential_mode" env:"CREDENTIAL_MODE"`
	// VaultPath Vault 路径：vault-secret 模式为 KV 路径，vault-dynamic 模式为角色名
	VaultPath string `json:"vault_path" yaml:"vault_path" toml:"vault_path" env:"VAULT_PATH"`
	// VaultUsernameField Vault 返回数据中的用户名字段名，默认 "username"
	VaultUsernameField string `json:"vault_username_field" yaml:"vault_username_field" toml:"vault_username_field" env:"VAULT_USERNAME_FIELD"`
	// VaultPasswordField Vault 返回数据中的密码字段名，默认 "password"
	VaultPasswordField string `json:"vault_password_field" yaml:"vault_password_field" toml:"vault_password_field" env:"VAULT_PASSWORD_FIELD"`
	// VaultAutoRenew 是否自动续期动态凭证，默认 true
	VaultAutoRenew bool `json:"vault_auto_renew" yaml:"vault_auto_renew" toml:"vault_auto_renew" env:"VAULT_AUTO_RENEW"`
	// VaultRenewThreshold 续期阈值（租约生命周期的百分比），默认 0.8 (80%)
	VaultRenewThreshold float64 `json:"vault_renew_threshold" yaml:"vault_renew_threshold" toml:"vault_renew_threshold" env:"VAULT_RENEW_THRESHOLD"`

	// GORM perform single create, update, delete operations in transactions by default to ensure database data integrity
	// You can disable it by setting `SkipDefaultTransaction` to true
	SkipDefaultTransaction bool `toml:"skip_default_transaction" json:"skip_default_transaction" yaml:"skip_default_transaction"  env:"SKIP_DEFALT_TRANSACTION"`
	// FullSaveAssociations full save associations
	FullSaveAssociations bool `toml:"full_save_associations" json:"full_save_associations" yaml:"full_save_associations"  env:"FULL_SAVE_ASSOCIATIONS"`
	// DryRun generate sql without execute
	DryRun bool `toml:"dry_run" json:"dry_run" yaml:"dry_run"  env:"DRY_RUN"`
	// PrepareStmt executes the given query in cached statement
	PrepareStmt bool `toml:"prepare_stmt" json:"prepare_stmt" yaml:"prepare_stmt"  env:"PREPARE_STMT"`
	// DisableAutomaticPing
	DisableAutomaticPing bool `toml:"disable_automatic_ping" json:"disable_automatic_ping" yaml:"disable_automatic_ping"  env:"DISABLE_AUTOMATIC_PING"`
	// DisableForeignKeyConstraintWhenMigrating
	DisableForeignKeyConstraintWhenMigrating bool `toml:"disable_foreign_key_constraint_when_migrating" json:"disable_foreign_key_constraint_when_migrating" yaml:"disable_foreign_key_constraint_when_migrating"  env:"DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING"`
	// IgnoreRelationshipsWhenMigrating
	IgnoreRelationshipsWhenMigrating bool `toml:"ignore_relationships_when_migrating" json:"ignore_relationships_when_migrating" yaml:"ignore_relationships_when_migrating"  env:"IGNORE_RELATIONSHIP_WHEN_MIGRATING"`
	// DisableNestedTransaction disable nested transaction
	DisableNestedTransaction bool `toml:"disable_nested_transaction" json:"disable_nested_transaction" yaml:"disable_nested_transaction"  env:"DISABLE_NESTED_TRANSACTION"`
	// AllowGlobalUpdate allow global update
	AllowGlobalUpdate bool `toml:"allow_global_update" json:"allow_global_update" yaml:"allow_global_update"  env:"ALL_GLOBAL_UPDATE"`
	// QueryFields executes the SQL query with all fields of the table
	QueryFields bool `toml:"query_fields" json:"query_fields" yaml:"query_fields"  env:"QUERY_FIELDS"`
	// CreateBatchSize default create batch size
	CreateBatchSize int `toml:"create_batch_size" json:"create_batch_size" yaml:"create_batch_size"  env:"CREATE_BATCH_SIZE"`
	// TranslateError enabling error translation
	TranslateError bool `toml:"translate_error" json:"translate_error" yaml:"translate_error"  env:"TRANSLATE_ERROR"`

	db  *gorm.DB
	log *zerolog.Logger

	// Vault 动态凭证内部状态
	leaseID       string        // vault-dynamic 模式的 lease ID
	leaseDuration int64         // 租约时长（秒）
	stopRenewal   chan struct{} // 停止续期信号
}

// toml，配置的名字datasource
func (m *dataSource) Name() string {
	return AppName
}

func (i *dataSource) Priority() int {
	return 999
}

func (m *dataSource) Init() error {
	m.log = log.Sub(m.Name())

	// 初始化默认值
	if err := m.setDefaults(); err != nil {
		return err
	}

	// 根据凭证模式加载凭证
	if err := m.loadCredentials(); err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}

	db, err := gorm.Open(m.Dialector(), &gorm.Config{
		SkipDefaultTransaction:                   m.SkipDefaultTransaction,
		FullSaveAssociations:                     m.FullSaveAssociations,
		DryRun:                                   m.DryRun,
		PrepareStmt:                              m.PrepareStmt,
		DisableAutomaticPing:                     m.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: m.DisableForeignKeyConstraintWhenMigrating,
		IgnoreRelationshipsWhenMigrating:         m.IgnoreRelationshipsWhenMigrating,
		DisableNestedTransaction:                 m.DisableNestedTransaction,
		AllowGlobalUpdate:                        m.AllowGlobalUpdate,
		Logger:                                   newGormCustomLogger(m.log),
	})
	if err != nil {
		return err
	}

	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable gorm trace")
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			return err
		}
	}

	if m.Debug {
		db = db.Debug()
	}

	m.db = db

	// 启动凭证续期（仅 vault-dynamic 模式）
	if m.CredentialMode.NeedsRenewal() && m.VaultAutoRenew {
		go m.startCredentialRenewal()
	}

	return nil
}

func (m *dataSource) Dialector() gorm.Dialector {
	switch m.Provider {
	case PROVIDER_POSTGRES:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			m.Host,
			m.Username,
			m.Password,
			m.DB,
			m.Port,
		)
		return postgres.Open(dsn)
	case PROVIDER_SQLITE:
		return sqlite.Open(m.DB)
	default:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.Username,
			m.Password,
			m.Host,
			m.Port,
			m.DB,
		)
		return mysql.Open(dsn)
	}
}

// 关闭数据库连接
func (m *dataSource) Close(ctx context.Context) error {
	// 停止续期
	if m.stopRenewal != nil {
		close(m.stopRenewal)
	}

	// 撤销动态凭证租约
	if m.CredentialMode == CREDENTIAL_MODE_VAULT_DYNAMIC && m.leaseID != "" {
		vaultClient := vault.Client()
		if vaultClient != nil {
			if _, err := vaultClient.System.LeasesRevokeLease(ctx, schema.LeasesRevokeLeaseRequest{LeaseId: m.leaseID}); err != nil {
				m.log.Error().Err(err).Msgf("failed to revoke lease %s", m.leaseID)
			} else {
				m.log.Info().Msgf("revoked vault lease %s", m.leaseID)
			}
		}
	}

	// 关闭数据库连接
	if m.db == nil {
		return nil
	}

	d, err := m.db.DB()
	if err != nil {
		m.log.Error().Msgf("获取db error, %s", err)
		return err
	}
	if err := d.Close(); err != nil {
		m.log.Error().Msgf("close db error, %s", err)
		return err
	}
	return nil
}

func (i *dataSource) Version() string {
	return ""
}
func (i *dataSource) AllowOverwrite() bool {
	return false
}
