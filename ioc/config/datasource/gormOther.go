package datasource

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault-client-go/schema"
	"github.com/kade-chen/library/ioc/config/trace"
	"github.com/kade-chen/library/ioc/config/vault"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
)

// setDefaults 设置默认值
func (m *dataSource) setDefaults() error {
	// Vault 字段默认值
	if m.VaultUsernameField == "" {
		m.VaultUsernameField = "username"
	}
	if m.VaultPasswordField == "" {
		m.VaultPasswordField = "password"
	}
	if m.VaultRenewThreshold == 0 {
		m.VaultRenewThreshold = 0.8
	}
	if m.VaultRenewThreshold < 0.5 || m.VaultRenewThreshold > 0.95 {
		return fmt.Errorf("vault_renew_threshold must be between 0.5 and 0.95, got %.2f", m.VaultRenewThreshold)
	}

	// 初始化续期停止信号
	if m.stopRenewal == nil {
		m.stopRenewal = make(chan struct{})
	}

	// 处理空字符串，兼容旧配置
	if m.CredentialMode == "" {
		m.CredentialMode = CREDENTIAL_MODE_STATIC
	}

	return nil
}

// loadCredentials 根据凭证模式加载凭证
func (m *dataSource) loadCredentials() error {
	switch m.CredentialMode {
	case CREDENTIAL_MODE_STATIC:
		return m.loadStaticCredentials()

	case CREDENTIAL_MODE_VAULT_SECRET:
		return m.loadVaultSecretCredentials()

	case CREDENTIAL_MODE_VAULT_DYNAMIC:
		return m.loadVaultDynamicCredentials()

	default:
		return fmt.Errorf("unsupported credential mode: %s", m.CredentialMode)
	}
}

// loadStaticCredentials 加载静态凭证
func (m *dataSource) loadStaticCredentials() error {
	if m.Username == "" {
		m.log.Warn().Msg("username is empty for static credential mode")
	}
	if m.Password == "" {
		m.log.Warn().Msg("password is empty for static credential mode")
	}
	m.log.Info().Msg("using static credentials from config file")
	return nil
}

// loadVaultSecretCredentials 从 Vault KV 加载静态凭证
func (m *dataSource) loadVaultSecretCredentials() error {
	if m.VaultPath == "" {
		return fmt.Errorf("vault_path is required for vault-secret mode")
	}

	vaultClient := vault.Client()
	if vaultClient == nil {
		return fmt.Errorf("vault client not initialized, please ensure vault config is enabled")
	}

	ctx := context.Background()
	resp, err := vault.ReadSecret(ctx, m.VaultPath)
	if err != nil {
		return fmt.Errorf("read vault secret from %s: %w", m.VaultPath, err)
	}

	// 提取 username 和 password
	username, ok := resp.Data.Data[m.VaultUsernameField].(string)
	if !ok {
		return fmt.Errorf("field '%s' not found in vault secret at %s", m.VaultUsernameField, m.VaultPath)
	}

	password, ok := resp.Data.Data[m.VaultPasswordField].(string)
	if !ok {
		return fmt.Errorf("field '%s' not found in vault secret at %s", m.VaultPasswordField, m.VaultPath)
	}

	m.Username = username
	m.Password = password

	m.log.Info().Msgf("loaded credentials from vault KV (path=%s)", m.VaultPath)
	return nil
}

// loadVaultDynamicCredentials 从 Vault Database 引擎生成动态凭证
func (m *dataSource) loadVaultDynamicCredentials() error {
	if m.VaultPath == "" {
		return fmt.Errorf("vault_path (role name) is required for vault-dynamic mode")
	}

	vaultClient := vault.Client()
	if vaultClient == nil {
		return fmt.Errorf("vault client not initialized, please ensure vault config is enabled")
	}

	ctx := context.Background()
	resp, err := vault.GenerateDatabaseCredentials(ctx, m.VaultPath)
	if err != nil {
		return fmt.Errorf("generate vault database credentials (role=%s): %w", m.VaultPath, err)
	}

	// 提取凭证
	username, ok := resp.Data["username"].(string)
	if !ok {
		return fmt.Errorf("username not found in vault database credentials response")
	}

	password, ok := resp.Data["password"].(string)
	if !ok {
		return fmt.Errorf("password not found in vault database credentials response")
	}

	m.Username = username
	m.Password = password

	// 提取租约信息
	m.leaseID = resp.LeaseID
	m.leaseDuration = int64(resp.LeaseDuration)

	m.log.Info().Msgf("generated vault dynamic credentials (role=%s, lease_id=%s, ttl=%ds)",
		m.VaultPath, m.leaseID, m.leaseDuration)

	return nil
}

// startCredentialRenewal 启动凭证自动续期（仅用于 vault-dynamic 模式）
func (m *dataSource) startCredentialRenewal() {
	if m.leaseDuration <= 0 {
		m.log.Warn().Msg("invalid lease duration, credential renewal disabled")
		return
	}

	// 在租约阈值时续期
	renewInterval := time.Duration(float64(m.leaseDuration)*m.VaultRenewThreshold) * time.Second

	m.log.Info().Msgf("credential auto-renewal enabled (interval=%s, threshold=%.0f%%)",
		renewInterval, m.VaultRenewThreshold*100)

	ticker := time.NewTicker(renewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := m.renewCredentials(); err != nil {
				m.log.Error().Err(err).Msg("failed to renew credentials")
				// 续期失败，尝试重新生成凭证
				if err := m.regenerateCredentials(); err != nil {
					m.log.Error().Err(err).Msg("failed to regenerate credentials - database connection may fail")
				}
			}

		case <-m.stopRenewal:
			m.log.Info().Msg("credential renewal stopped")
			return
		}
	}
}

// renewCredentials 续期现有凭证租约
func (m *dataSource) renewCredentials() error {
	ctx := context.Background()
	client := vault.Client()

	resp, err := client.System.LeasesRenewLease(ctx, schema.LeasesRenewLeaseRequest{
		LeaseId:   m.leaseID,
		Increment: fmt.Sprintf("%ds", m.leaseDuration),
	})
	if err != nil {
		return fmt.Errorf("renew lease %s: %w", m.leaseID, err)
	}

	m.leaseDuration = int64(resp.LeaseDuration)
	m.log.Info().Msgf("credentials renewed (lease_id=%s, new_ttl=%ds)", m.leaseID, m.leaseDuration)

	return nil
}

// regenerateCredentials 重新生成凭证并重连数据库
func (m *dataSource) regenerateCredentials() error {
	m.log.Warn().Msg("attempting to regenerate credentials")

	// 重新生成凭证
	if err := m.loadVaultDynamicCredentials(); err != nil {
		return fmt.Errorf("regenerate credentials: %w", err)
	}

	// 关闭旧连接
	if m.db != nil {
		if sqlDB, err := m.db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	// 创建新连接
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
		return fmt.Errorf("reconnect database with new credentials: %w", err)
	}

	if trace.Get().Enable && m.Trace {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			return fmt.Errorf("enable trace on new connection: %w", err)
		}
	}

	if m.Debug {
		db = db.Debug()
	}

	m.db = db
	m.log.Info().Msg("database reconnected with regenerated credentials")

	return nil
}
