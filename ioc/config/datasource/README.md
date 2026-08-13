# 关系型数据库（datasource）

基于 [GORM]() 的关系型数据库配置模块，集成到 mcube IoC 容器中，支持 MySQL、PostgreSQL、SQLite 三种数据库驱动，以及静态凭证、Vault KV 静态凭证、Vault 动态凭证三种凭证加载模式。

## 目录

- [快速开始](#快速开始)
- [配置方式](#配置方式)
  - [环境变量](#环境变量)
  - [TOML 文件](#toml-文件)
  - [配置优先级](#配置优先级)
- [数据库驱动](#数据库驱动)
  - [MySQL](#mysql)
  - [PostgreSQL](#postgresql)
  - [SQLite](#sqlite)
- [凭证模式](#凭证模式)
  - [静态凭证（static）](#静态凭证static)
  - [Vault KV 静态凭证（vault-secret）](#vault-kv-静态凭证vault-secret)
  - [Vault 动态凭证（vault-dynamic）](#vault-动态凭证vault-dynamic)
- [事务管理](#事务管理)
- [GORM 高级配置](#gorm-高级配置)
- [链路追踪](#链路追踪)
- [完整配置参数](#完整配置参数)

---

## 快速开始

**1. 导入包（触发自动注册）**

```go
import (
    "github.com/kade-chen/library/ioc"
    "github.com/kade-chen/library/ioc/config/datasource"
    _ "github.com/kade-chen/library/ioc/config/datasource" // 仅注册，不使用
)
```

**2. 通过环境变量提供配置，然后初始化**

```go
func main() {
    // 方式一：开发环境，从配置文件加载
    ioc.DevelopmentSetupWithPath("etc/application.toml")

    // 方式二：生产环境，从环境变量加载
    ioc.ConfigIocObject(ioc.NewLoadConfigRequest())

    // 获取 *gorm.DB 对象
    db := datasource.DB()
    fmt.Println(db)
}
```

**3. 最简 TOML 配置（`etc/application.toml`）**

```toml
[datasource]
  host     = "127.0.0.1"
  port     = 3306
  database = "myapp"
  username = "root"
  password = "123456"
```

---

## 配置方式

### 环境变量

所有配置项均可通过 `DATASOURCE_` 前缀的环境变量设置。

```bash
export DATASOURCE_HOST=127.0.0.1
export DATASOURCE_PORT=3306
export DATASOURCE_DB=myapp
export DATASOURCE_USERNAME=root
export DATASOURCE_PASSWORD=123456
export DATASOURCE_PROVIDER=mysql   # mysql | postgres | sqlite
export DATASOURCE_DEBUG=false
export DATASOURCE_TRACE=true
```

凭证模式相关：

```bash
export DATASOURCE_CREDENTIAL_MODE=static   # static | vault-secret | vault-dynamic
export DATASOURCE_VAULT_PATH=secret/data/myapp/db
export DATASOURCE_VAULT_USERNAME_FIELD=username
export DATASOURCE_VAULT_PASSWORD_FIELD=password
export DATASOURCE_VAULT_AUTO_RENEW=true
export DATASOURCE_VAULT_RENEW_THRESHOLD=0.8
```

### TOML 文件

```toml
[datasource]
  provider = "mysql"           # 数据库驱动，默认 mysql
  host     = "127.0.0.1"
  port     = 3306
  database = "myapp"
  username = "root"
  password = "123456"
  debug    = false             # 打印所有 SQL
  trace    = true              # 启用 OpenTelemetry 链路追踪

  # 凭证模式（可选，默认 static）
  credential_mode      = "static"
  vault_path           = ""
  vault_username_field = "username"
  vault_password_field = "password"
  vault_auto_renew     = true
  vault_renew_threshold = 0.8

  # GORM 高级选项（可选）
  skip_default_transaction = false
  prepare_stmt             = true
  dry_run                  = false
```

### 配置优先级

当同时使用文件配置和环境变量时，**环境变量优先级更高**，会覆盖文件中相同的配置项。

```
环境变量 > TOML/YAML 文件 > 代码默认值
```

---

## 数据库驱动

通过 `provider` 字段（或 `DATASOURCE_PROVIDER` 环境变量）切换驱动，默认为 `mysql`。

### MySQL

```toml
[datasource]
  provider = "mysql"
  host     = "127.0.0.1"
  port     = 3306
  database = "myapp"
  username = "root"
  password = "123456"
```

DSN 格式：`{username}:{password}@tcp({host}:{port})/{database}?charset=utf8mb4&parseTime=True&loc=Local`

> 依赖：`gorm.io/driver/mysql`（已内置）

### PostgreSQL

```toml
[datasource]
  provider = "postgres"
  host     = "127.0.0.1"
  port     = 5432
  database = "myapp"
  username = "postgres"
  password = "123456"
```

DSN 格式：`host={host} user={username} password={password} dbname={database} port={port} sslmode=disable TimeZone=Asia/Shanghai`

> 依赖：`gorm.io/driver/postgres`（已内置）

### SQLite

SQLite 模式下 `database` 字段填写数据库文件路径，不需要填写 `host`、`port`、`username`、`password`。

```toml
[datasource]
  provider = "sqlite"
  database = "data/myapp.db"   # SQLite 文件路径
```

> 依赖：`github.com/glebarez/sqlite`（纯 Go 实现，无需 CGO）

---

## 凭证模式

通过 `credential_mode` 字段（或 `DATASOURCE_CREDENTIAL_MODE` 环境变量）控制数据库凭证的来源，支持三种模式。

### 静态凭证（static）

**默认模式**。直接从配置文件或环境变量中读取 `username` 和 `password`。

```toml
[datasource]
  credential_mode = "static"   # 默认，可省略
  username = "root"
  password = "123456"
```

- ✅ 简单直接，适合开发和小规模部署
- ❌ 密码明文存储在配置文件中，安全性较低
- ❌ 轮换密码需要重启服务

---

### Vault KV 静态凭证（vault-secret）

从 HashiCorp Vault 的 **KV v2 引擎**中读取静态凭证。密码由 Vault 统一管理，应用不再直接持有密码。

**配置示例：**

```toml
[datasource]
  credential_mode      = "vault-secret"
  host                 = "127.0.0.1"
  port                 = 3306
  database             = "myapp"
  vault_path           = "myapp/db"        # KV 路径（相对于挂载点）
  vault_username_field = "username"         # Vault 返回数据中的用户名字段，默认 "username"
  vault_password_field = "password"         # Vault 返回数据中的密码字段，默认 "password"

[vault]
  address     = "http://127.0.0.1:8200"
  auth_method = "token"
  token       = "hvs.xxxxxx"
```

**Vault 侧配置（参考）：**

```bash
# 1. 启动 Vault 开发服务
docker run -itd --name=vault \
  --cap-add=IPC_LOCK \
  -p 8200:8200 \
  -e 'VAULT_DEV_ROOT_TOKEN_ID=myroot' \
  -e 'VAULT_DEV_LISTEN_ADDRESS=0.0.0.0:8200' \
  hashicorp/vault:latest

# 2. 启用 KV v2 引擎（开发模式已自动启用，默认挂载点为 secret）
docker exec \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  -e VAULT_TOKEN='myroot' \
  vault vault secrets enable -path=secret kv-v2

# 3. 写入数据库凭证
docker exec \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  -e VAULT_TOKEN='myroot' \
  vault vault kv put secret/myapp/db \
    username=root \
    password=123456

# 4. 验证读取
docker exec \
  -e VAULT_ADDR='http://127.0.0.1:8200' \
  -e VAULT_TOKEN='myroot' \
  vault vault kv get secret/myapp/db
```

> `vault_path` 填写的是相对于挂载点的路径。默认挂载点为 `secret`（对应 `kv_mount_path` 配置项），
> 实际请求路径为 `secret/data/myapp/db`，此处只需填 `myapp/db`。

**工作流程：**

```
应用启动
  └─→ 连接 Vault
  └─→ 读取 KV 路径中的 username / password
  └─→ 使用读取到的凭证连接数据库
```

- ✅ 密码不出现在应用配置文件中
- ✅ 可以在 Vault 中集中管理和轮换密码
- ❌ 密码仍然是静态的，轮换后需要重启服务

---

### Vault 动态凭证（vault-dynamic）

通过 Vault 的 **Database 引擎**为每次应用启动**动态生成**一个临时数据库账号，凭证有租约（TTL），到期自动续期或重新生成。

**配置示例：**

```toml
[datasource]
  credential_mode       = "vault-dynamic"
  host                  = "127.0.0.1"
  port                  = 3306
  database              = "myapp"
  vault_path            = "my-role"          # Vault Database 引擎角色名
  vault_auto_renew      = true               # 自动续期，默认 true
  vault_renew_threshold = 0.8                # 在租约剩余 20% 时续期，默认 0.8

[vault]
  address     = "http://127.0.0.1:8200"
  auth_method = "token"
  token       = "hvs.xxxxxx"
```

**Vault 侧配置（参考）：**

```bash
# 1. 启用 Database 引擎
vault secrets enable database

# 2. 配置数据库连接
vault write database/config/myapp \
    plugin_name=mysql-database-plugin \
    connection_url="{{username}}:{{password}}@tcp(127.0.0.1:3306)/" \
    allowed_roles="my-role" \
    username="vault_admin" \
    password="vault_admin_pass"

# 3. 创建角色
vault write database/roles/my-role \
    db_name=myapp \
    creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT,INSERT,UPDATE,DELETE ON myapp.* TO '{{name}}'@'%';" \
    default_ttl="1h" \
    max_ttl="24h"
```

**凭证生命周期：**

```
应用启动
  └─→ 调用 Vault Database API 生成临时凭证（username + password + lease_id）
  └─→ 使用临时凭证连接数据库
  └─→ 后台 goroutine 监控租约，在 TTL × 80% 时自动续期
        └─→ 续期失败 → 重新生成新凭证 → 重建数据库连接
应用关闭
  └─→ 调用 Vault API 立即撤销租约
```

**续期机制：**

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `vault_auto_renew` | `true` | 是否启用自动续期 |
| `vault_renew_threshold` | `0.8` | 续期触发时机（租约生命周期的百分比，范围 0.5~0.95） |

- ✅ 凭证自动过期，大幅降低泄露风险
- ✅ 生产环境最佳实践，符合零信任安全模型
- ✅ 续期失败时自动重试，无需人工干预
- ❌ 依赖 Vault 服务，增加基础设施复杂度

---

## 事务管理

模块内置了基于 `context.Context` 的事务传播机制，无需在业务函数间显式传递 `*gorm.DB`。

**开启事务并传播：**

```go
err := datasource.DB().Transaction(func(tx *gorm.DB) error {
    // 将事务绑定到 context
    txCtx := datasource.WithTransactionCtx(ctx, tx)

    // 调用其他业务方法，事务自动透传
    if err := createOrder(txCtx); err != nil {
        return err  // 返回 error 则自动回滚
    }
    if err := deductStock(txCtx); err != nil {
        return err
    }
    return nil  // 返回 nil 则自动提交
})
```

**业务方法使用事务：**

```go
func createOrder(ctx context.Context) error {
    // DBFromCtx：如果 ctx 中有事务则使用事务，否则使用普通 DB
    db := datasource.DBFromCtx(ctx)
    return db.Create(&Order{...}).Error
}

func deductStock(ctx context.Context) error {
    db := datasource.DBFromCtx(ctx)
    return db.Model(&Stock{}).Where("id = ?", 1).Update("count", gorm.Expr("count - 1")).Error
}
```

**API 说明：**

| 函数 | 说明 |
|------|------|
| `datasource.DB()` | 获取全局 `*gorm.DB` 实例 |
| `datasource.DBFromCtx(ctx)` | 从 context 中取事务，取不到则返回普通 DB |
| `datasource.WithTransactionCtx(ctx, tx)` | 将事务绑定到 context |

---

## GORM 高级配置

以下参数对应 [gorm.Config](https://gorm.io/docs/gorm_config.html)，可按需调整。

```toml
[datasource]
  # 是否跳过默认事务（对 Create/Update/Delete 的单操作事务包裹）
  # 开启后性能更高，但失去单操作原子性保障，默认 false
  skip_default_transaction = false

  # 是否使用预编译 SQL 缓存，减少 SQL 解析开销，默认 true
  prepare_stmt = true

  # 是否开启 SQL 打印模式（仅打印，不执行），用于调试，默认 false
  dry_run = false

  # 是否保存全量关联关系，默认 false
  full_save_associations = false

  # 是否禁用 Ping，默认 false
  disable_automatic_ping = false

  # 迁移时是否跳过外键约束，默认 false
  disable_foreign_key_constraint_when_migrating = false

  # 迁移时是否跳过关系，默认 false
  ignore_relationships_when_migrating = false

  # 是否禁用嵌套事务，默认 false
  disable_nested_transaction = false

  # 是否允许全局更新（不带 WHERE 条件），默认 false（强烈不建议开启）
  allow_global_update = false

  # 查询时是否使用 SELECT * 展开所有字段，默认 false
  query_fields = false

  # 批量创建时的默认批次大小，0 表示不限制
  create_batch_size = 0

  # 是否将 GORM 错误翻译为人类可读格式，默认 false
  translate_error = false

  # 是否打印所有 SQL（等同于 db.Debug()），默认 false
  debug = false
```

---

## 链路追踪

当 `trace` 配置项为 `true` 且 mcube 的 `trace` 模块也启用了 OpenTelemetry 时，所有数据库操作会自动创建 Span，记录：

- SQL 语句
- 执行耗时
- 数据库连接信息

```toml
[datasource]
  trace = true    # 开启数据库级 trace，默认 true

[trace]
  enable   = true
  endpoint = "http://jaeger:14268/api/traces"
```

> trace 集成基于 [otelgorm](https://github.com/uptrace/opentelemetry-go-extra/tree/main/otelgorm)，与 mcube 的 `trace` 模块共享同一 TracerProvider，业务请求的 Span 会自动成为数据库 Span 的父节点。

---

## 完整配置参数

| 参数 | 环境变量 | 类型 | 默认值 | 说明 |
|------|----------|------|--------|------|
| `provider` | `DATASOURCE_PROVIDER` | string | `mysql` | 驱动：`mysql` / `postgres` / `sqlite` |
| `host` | `DATASOURCE_HOST` | string | `127.0.0.1` | 数据库地址 |
| `port` | `DATASOURCE_PORT` | int | `3306` | 数据库端口 |
| `database` | `DATASOURCE_DB` | string | 应用名 | 数据库名（SQLite 为文件路径） |
| `username` | `DATASOURCE_USERNAME` | string | — | 用户名 |
| `password` | `DATASOURCE_PASSWORD` | string | — | 密码 |
| `auto_migrate` | `DATASOURCE_AUTO_MIGRATE` | bool | `false` | 自动执行 GORM AutoMigrate（需业务代码配合） |
| `debug` | `DATASOURCE_DEBUG` | bool | `false` | 打印所有 SQL |
| `trace` | `DATASOURCE_TRACE` | bool | `true` | 启用 OpenTelemetry 链路追踪 |
| `credential_mode` | `DATASOURCE_CREDENTIAL_MODE` | string | `static` | 凭证模式：`static` / `vault-secret` / `vault-dynamic` |
| `vault_path` | `DATASOURCE_VAULT_PATH` | string | — | Vault 路径（KV 路径或角色名） |
| `vault_username_field` | `DATASOURCE_VAULT_USERNAME_FIELD` | string | `username` | Vault 返回数据中的用户名字段 |
| `vault_password_field` | `DATASOURCE_VAULT_PASSWORD_FIELD` | string | `password` | Vault 返回数据中的密码字段 |
| `vault_auto_renew` | `DATASOURCE_VAULT_AUTO_RENEW` | bool | `true` | 动态凭证自动续期 |
| `vault_renew_threshold` | `DATASOURCE_VAULT_RENEW_THRESHOLD` | float64 | `0.8` | 续期触发阈值（0.5~0.95） |
| `skip_default_transaction` | `DATASOURCE_SKIP_DEFALT_TRANSACTION` | bool | `false` | 跳过单操作默认事务 |
| `prepare_stmt` | `DATASOURCE_PREPARE_STMT` | bool | `true` | 预编译 SQL 缓存 |
| `dry_run` | `DATASOURCE_DRY_RUN` | bool | `false` | 仅生成 SQL，不执行 |
| `full_save_associations` | `DATASOURCE_FULL_SAVE_ASSOCIATIONS` | bool | `false` | 保存全量关联 |
| `disable_automatic_ping` | `DATASOURCE_DISABLE_AUTOMATIC_PING` | bool | `false` | 禁用自动 Ping |
| `disable_foreign_key_constraint_when_migrating` | `DATASOURCE_DISABLE_FOREIGN_KEY_CONSTRAINT_WHEN_MIGRATING` | bool | `false` | 迁移时跳过外键约束 |
| `ignore_relationships_when_migrating` | `DATASOURCE_IGNORE_RELATIONSHIP_WHEN_MIGRATING` | bool | `false` | 迁移时跳过关系 |
| `disable_nested_transaction` | `DATASOURCE_DISABLE_NESTED_TRANSACTION` | bool | `false` | 禁用嵌套事务 |
| `allow_global_update` | `DATASOURCE_ALL_GLOBAL_UPDATE` | bool | `false` | 允许无 WHERE 的全局更新 |
| `query_fields` | `DATASOURCE_QUERY_FIELDS` | bool | `false` | 查询时展开所有字段 |
| `create_batch_size` | `DATASOURCE_CREATE_BATCH_SIZE` | int | `0` | 默认批量创建大小 |
| `translate_error` | `DATASOURCE_TRANSLATE_ERROR` | bool | `false` | 翻译 GORM 错误信息 |