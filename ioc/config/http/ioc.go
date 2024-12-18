package http

import (
	"net/http"
	"time"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Http{
	Host:                    "127.0.0.1",
	Port:                    8080,
	PathPrefix:              "api",
	ReadHeaderTimeoutSecond: 30,
	ReadTimeoutSecond:       60,
	WriteTimeoutSecond:      60,
	IdleTimeoutSecond:       600,
	MaxHeaderSize:           "16kb",
	EnableTrace:             true,
}

type Http struct {
	ioc.ObjectImpl

	// 是否开启HTTP Server, 默认会根据是否有注册得有API对象来自动开启
	Enable *bool `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	// HTTP服务Host
	Host string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	// HTTP服务端口
	Port int `json:"port" yaml:"port" toml:"port" env:"PORT"`
	// API接口前缀
	PathPrefix string `json:"path_prefix" yaml:"path_prefix" toml:"path_prefix" env:"PATH_PREFIX"`

	// HTTP服务器参数
	// HTTP Header读取超时时间
	ReadHeaderTimeoutSecond int `json:"read_header_timeout" yaml:"read_header_timeout" toml:"read_header_timeout" env:"READ_HEADER_TIMEOUT"`
	// 读取HTTP整个请求时的参数
	ReadTimeoutSecond int `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" env:"READ_TIMEOUT"`
	// 响应超时时间
	WriteTimeoutSecond int `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" env:"WRITE_TIMEOUT"`
	// 启用了KeepAlive时 复用TCP链接的超时时间
	IdleTimeoutSecond int `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" env:"IDLE_TIMEOUT"`
	// header最大大小
	MaxHeaderSize string `json:"max_header_size" yaml:"max_header_size" toml:"max_header_size" env:"MAX_HEADER_SIZE"`

	// 开启Trace
	EnableTrace bool `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace" env:"ENABLE_TRACE"`
	// 开启Trace
	EnableDebug bool `toml:"enable_debug" json:"enable_debug" yaml:"enable_debug" env:"ENABLE_DEBUG"`

	// 解析后的数据
	maxHeaderBytes uint64
	log            *zerolog.Logger
	router         http.Handler
	server         *http.Server
}

func (h *Http) Name() string {
	return AppName
}

func (i *Http) Priority() int {
	return -9998
}

// 配置数据解析
func (h *Http) Init() error {
	h.log = log.Sub("http")

	mhz, err := humanize.ParseBytes(h.MaxHeaderSize)
	if err != nil {
		return err
	}
	h.maxHeaderBytes = mhz

	h.server = &http.Server{
		ReadHeaderTimeout: time.Duration(h.ReadHeaderTimeoutSecond) * time.Second,
		ReadTimeout:       time.Duration(h.ReadTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(h.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(h.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    int(h.maxHeaderBytes),
		Addr:              h.Addr(),
		Handler:           h.router,
	}
	return nil
}
