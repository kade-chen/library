package redis

import (
	"context"
	"time"

	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/kade-chen/library/ioc/config/trace"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Redis{
	DB:        0,
	Endpoints: []string{"127.0.0.1:6379"},
	Trace:     false,
}

type Redis struct {
	// IOC基础对象（依赖注入用）
	ioc.ObjectImpl
	// Endpoints Redis地址列表（支持单机/集群）
	// 单机: ["127.0.0.1:6379"]
	// 集群: ["10.0.0.1:6379","10.0.0.2:6379"]
	Endpoints []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"ENDPOINTS" envSeparator:","`
	// DB 使用的Redis数据库编号（默认0）
	DB int `toml:"db" json:"db" yaml:"db" env:"DB"`
	// UserName Redis ACL用户名（Redis 6+）
	UserName string `toml:"username" json:"username" yaml:"username" env:"USERNAME"`
	// Password Redis密码
	Password string `toml:"password" json:"password" yaml:"password" env:"PASSWORD"`
	// ================= 性能参数 =================
	// PoolSize 连接池大小（最重要参数）
	// 决定并发能力，建议：100~300（根据QPS调整）
	PoolSize int `toml:"pool_size" json:"pool_size" yaml:"pool_size" env:"POOL_SIZE"`
	// MinIdleConns 最小空闲连接数
	// 作用：减少频繁建连，提高响应速度
	MinIdleConns int `toml:"min_idle_conns" json:"min_idle_conns" yaml:"min_idle_conns" env:"MIN_IDLE_CONNS"`
	// ================= 超时控制（单位：毫秒） =================
	// DialTimeout 建立连接超时时间
	// Redis不可达时快速失败
	DialTimeout int `toml:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout" env:"DIAL_TIMEOUT"`
	// ReadTimeout 读取超时
	// 防止Redis卡顿导致请求堆积
	ReadTimeout int `toml:"read_timeout" json:"read_timeout" yaml:"read_timeout" env:"READ_TIMEOUT"`
	// WriteTimeout 写入超时
	// 防止写请求阻塞
	WriteTimeout int `toml:"write_timeout" json:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	// PoolTimeout 从连接池获取连接的超时时间
	// 防止连接耗尽导致阻塞
	PoolTimeout int `toml:"pool_timeout" json:"pool_timeout" yaml:"pool_timeout" env:"POOL_TIMEOUT"`
	// ================= 重试机制 =================
	// MaxRetries 最大重试次数
	// 建议：1~2，过大会影响延迟
	MaxRetries int `toml:"max_retries" json:"max_retries" yaml:"max_retries" env:"MAX_RETRIES"`
	// ================= 连接生命周期 =================
	// ConnMaxIdleTime 空闲连接最大存活时间（秒）
	// 超过会被关闭，释放资源
	ConnMaxIdleTime int `toml:"conn_max_idle_time" json:"conn_max_idle_time" yaml:"conn_max_idle_time" env:"CONN_MAX_IDLE_TIME"`
	// ConnMaxLifetime 连接最大生命周期（秒）
	// 防止连接长期使用导致不稳定
	ConnMaxLifetime int `toml:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime" env:"CONN_MAX_LIFETIME"`
	// ================= 可观测 =================
	// Trace 是否开启链路追踪（Tracing）
	Trace bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`
	// Metric 是否开启指标统计（QPS、延迟等）
	Metric bool `toml:"metric" json:"metric" yaml:"metric" env:"METRIC"`
	// ================= 内部对象 =================
	// client Redis客户端（支持单机/集群/哨兵）
	client redis.UniversalClient
	// log 日志组件
	log *zerolog.Logger
}

func (m *Redis) Name() string {
	return AppName
}

func (i *Redis) Priority() int {
	return 697
}

// https://opentelemetry.io/ecosystem/registry/?s=redis&component=&language=go
// https://github.com/redis/go-redis/tree/master/extra/redisotel
func (m *Redis) Init() error {
	m.log = log.Sub(m.Name())

	// 创建UniversalClient（支持单机/集群/哨兵）
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		// Addrs Redis地址
		Addrs: m.Endpoints,
		// DB 数据库编号
		DB: m.DB,
		// Username ACL用户名
		Username: m.UserName,
		// Password 密码
		Password: m.Password,
		// PoolSize 最大连接数（并发能力）
		PoolSize: defaultIfZero(m.PoolSize, 200),
		// MinIdleConns 最小空闲连接
		MinIdleConns: defaultIfZero(m.MinIdleConns, 20),
		// DialTimeout 建连超时（ms）
		DialTimeout: time.Duration(defaultIfZero(m.DialTimeout, 1000)) * time.Millisecond,
		// ReadTimeout 读超时（ms）
		ReadTimeout: time.Duration(defaultIfZero(m.ReadTimeout, 500)) * time.Millisecond,
		// WriteTimeout 写超时（ms）
		WriteTimeout: time.Duration(defaultIfZero(m.WriteTimeout, 500)) * time.Millisecond,
		// PoolTimeout 取连接等待（ms）
		PoolTimeout: time.Duration(defaultIfZero(m.PoolTimeout, 1000)) * time.Millisecond,
		// MaxRetries 最大重试次数
		MaxRetries: defaultIfZero(m.MaxRetries, 1),
		// ConnMaxIdleTime 空闲连接存活（s）
		ConnMaxIdleTime: time.Duration(defaultIfZero(m.ConnMaxIdleTime, 300)) * time.Second,
		// ConnMaxLifetime 连接最大生命周期（s）
		ConnMaxLifetime: time.Duration(defaultIfZero(m.ConnMaxLifetime, 1800)) * time.Second,
	})

	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable redis trace")
		if err := redisotel.InstrumentTracing(client); err != nil {
			return err
		}
	}

	if m.Metric {
		if err := redisotel.InstrumentMetrics(client); err != nil {
			return err
		}
	}

	m.client = client
	return nil
}

// 关闭数据库连接
func (m *Redis) Close(ctx context.Context) error {
	if m.client == nil {
		return exception.NewIocRegisterFailed("redis ")
	}

	err := m.client.Close()
	if err != nil {
		m.log.Error().Msgf("close redis client error, %s", err)
	}
	return nil
}

// defaultIfZero 如果配置为0，则使用默认值
func defaultIfZero(val int, def int) int {
	if val == 0 {
		return def
	}
	return val
}
