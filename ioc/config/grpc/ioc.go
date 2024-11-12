package grpc

import (
	"context"

	"github.com/kade-chen/library/grpc/middleware/recovery"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/kade-chen/library/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func init() {
	ioc.Config().Registry(&Grpc{
		Host:           "127.0.0.1",
		Port:           18080,
		EnableRecovery: true,
		EnableTrace:    true,
	})
}

type Grpc struct {
	ioc.ObjectImpl

	// 开启GRPC服务
	Enable *bool  `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Host   string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port   int    `json:"port" yaml:"port" toml:"port" env:"PORT"`

	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"KEY_FILE"`

	// 开启recovery恢复
	EnableRecovery bool `json:"enable_recovery" yaml:"enable_recovery" toml:"enable_recovery" env:"ENABLE_RECOVERY"`
	// 开启Trace
	EnableTrace bool `json:"enable_trace" yaml:"enable_trace" toml:"enable_trace" env:"ENABLE_TRACE"`

	// 解析后的数据
	interceptors []grpc.UnaryServerInterceptor
	svr          *grpc.Server
	log          *zerolog.Logger

	// 启动后执行
	PostStart func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
	// 关闭前执行
	PreStop func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
}

func (g *Grpc) Name() string {
	return AppName
}

func (i *Grpc) Priority() int {
	return 9997
}

func (g *Grpc) Init() error {
	g.log = log.Sub("grpc")
	g.svr = grpc.NewServer(g.ServerOpts()...)
	return nil
}

func (g *Grpc) ServerOpts() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	// 补充Trace选项
	if trace.Get().Enable && g.EnableTrace {
		// 初始化 OpenTelemetry gRPC 服务器处理器
		g.log.Info().Msg("enable mongodb trace")
		otelgrpc.NewServerHandler()
		// 添加跟踪处理器到选项列表
		//它用于创建 OpenTelemetry gRPC 服务器处理器。处理器是一个用于跟踪 gRPC 服务器请求和响应的组件。
		opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}
	// 补充中间件
	//这行代码将 gRPC 服务器的一组中间件 (拦截器) 添加到选项列表中
	opts = append(opts, grpc.ChainUnaryInterceptor(g.Interceptors()...))
	return opts
}

// Interceptors 返回 gRPC 服务器的一组拦截器
func (g *Grpc) Interceptors() (interceptors []grpc.UnaryServerInterceptor) {
	if g.EnableRecovery {
		//这行代码创建了一个恢复拦截器，并将其添加到拦截器列表中。恢复拦截器用于处理在 gRPC 服务器处理请求期间发生的 panic
		interceptors = append(interceptors, recovery.NewInterceptor(recovery.NewZeroLogRecoveryHandler()).
			UnaryServerInterceptor())
	}

	interceptors = append(interceptors, g.interceptors...)
	return
}
