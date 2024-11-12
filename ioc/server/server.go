package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"

	"github.com/kade-chen/library/ioc/config/grpc"
	"github.com/kade-chen/library/ioc/config/http" //里面用到了"github.com/kade-chen/library/ioc/config/trace" restful.go
	"github.com/rs/zerolog"
)

var DefaultConfig = ioc.NewLoadConfigRequest()

type Server struct {
	ioc.ObjectImpl
	setupHook func()

	http *http.Http
	grpc *grpc.Grpc

	ch     chan os.Signal
	log    *zerolog.Logger
	ctx    context.Context
	cancle context.CancelFunc
}

// req := ioc.NewLoadConfigRequest()
// req.ConfigFile.Enabled = true
// req.ConfigFile.Path = "/Users/kade.chen/wondorcloud/conf/application.toml"
func Run(ctx context.Context, lcr *ioc.LoadConfigRequest) error {
	return NewServer().Run(ctx, lcr)
}

func NewServer() *Server {
	return &Server{}
}

// func (a *Server) WithSetup(setup func()) *Server {
// 	a.setupHook = setup
// 	return a
// }

func (s *Server) Run(ctx context.Context, lcr *ioc.LoadConfigRequest) error {
	// 1.初始化ioc
	if lcr == nil {
		lcr = DefaultConfig
	}
	ioc.DevelopmentSetup(lcr)

	// 2.ioc setup 处理信号
	s.setup()

	s.log.Info().Msgf("loaded configs: %s", ioc.Config().List())
	s.log.Info().Msgf("loaded controllers: %s", ioc.Controller().List())
	s.log.Info().Msgf("loaded apis: %s", ioc.Api().List())

	if s.http.IsEnable() {
		go s.http.Start(ctx)
	}
	// //如果这里不给bool,指针默认给的是nil，就不能开启grpc
	// var grpc_bool bool
	if s.grpc.IsEnable() {
		go s.grpc.Start(ctx)
	}
	//等待取消信号
	s.waitSign()
	return nil
}

// 2.ioc setup 处理信号
func (s *Server) setup() {
	//处理信号
	s.ch = make(chan os.Signal, 1)

	//可以通过 signal.Notify 函数来注册一个或多个信号处理函数，以便在接收到指定的信号时执行相应的操作
	//它的作用是让程序能够捕获并处理操作系统发送的信号。
	signal.Notify(s.ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)

	//它的作用是创建一个新的上下文，并返回一个取消函数，当调用该函数时，会取消上下文。
	s.ctx, s.cancle = context.WithCancel(context.Background())

	s.http = http.Get()
	s.grpc = grpc.Get()
	s.log = log.Sub("server")
	if s.setupHook != nil {
		s.setupHook()
	}
}

// 3.等待取消信号
func (s *Server) waitSign() {
	defer s.cancle()

	for sg := range s.ch {
		switch v := sg.(type) {
		default:
			s.log.Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if s.grpc.IsEnable() {
				if err := s.grpc.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("grpc graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Info().Msg("grpc service stop complete")
				}
			}

			if s.http.IsEnable() {
				if err := s.http.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("http graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Info().Msgf("http service stop complete")
				}
			}
			return
		}
	}
}
