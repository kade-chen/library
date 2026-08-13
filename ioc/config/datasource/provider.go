package datasource

type PROVIDER string

const (
	PROVIDER_MYSQL    PROVIDER = "mysql"
	PROVIDER_POSTGRES PROVIDER = "postgres"
	PROVIDER_SQLITE   PROVIDER = "sqlite"
)

type CREDENTIAL_MODE string

const (
	CREDENTIAL_MODE_STATIC        CREDENTIAL_MODE = "static"        // 静态凭证（配置文件）
	CREDENTIAL_MODE_VAULT_SECRET  CREDENTIAL_MODE = "vault-secret"  // Vault KV 静态凭证
	CREDENTIAL_MODE_VAULT_DYNAMIC CREDENTIAL_MODE = "vault-dynamic" // Vault Database 动态凭证
)

// IsVaultMode 判断是否为 Vault 模式
func (m CREDENTIAL_MODE) IsVaultMode() bool {
	return m == CREDENTIAL_MODE_VAULT_SECRET || m == CREDENTIAL_MODE_VAULT_DYNAMIC
}

// NeedsRenewal 判断是否需要凭证续期
func (m CREDENTIAL_MODE) NeedsRenewal() bool {
	return m == CREDENTIAL_MODE_VAULT_DYNAMIC
}
