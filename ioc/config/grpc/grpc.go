package grpc

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

// -----------------------------------------------

type ServiceInfoCtxKey struct{}

func (g *Grpc) Start(ctx context.Context) {
	// 启动GRPC服务
	lis, err := net.Listen("tcp", g.Addr())
	if err != nil {

		g.log.Error().Msgf("listen grpc tcp conn error, %s", err)
		return
	}

	// 启动后勾子
	ctx = context.WithValue(ctx, ServiceInfoCtxKey{}, g.svr.GetServiceInfo())
	if g.PostStart != nil {
		if err := g.PostStart(ctx); err != nil {
			g.log.Error().Msg(err.Error())
			return
		}
	}

	g.log.Info().Msgf("GRPC 服务监听地址: http://%s", g.Addr())

	if err := g.svr.Serve(lis); err != nil {
		g.log.Error().Msg(err.Error())
	}
}

func (g *Grpc) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

// -----------------------------------------------

func (g *Grpc) Server() *grpc.Server {
	if g.svr == nil {
		panic("gprc server not initital")
	}
	return g.svr
}

// -----------------------------------------------

func (g *Grpc) Stop(ctx context.Context) error {
	// 停止之前的Hook
	if g.PreStop != nil {
		if err := g.PreStop(ctx); err != nil {
			return err
		}
	}

	g.svr.GracefulStop()
	return nil
}

// -----------------------------------------------

func (g *Grpc) IsEnable() bool {
	if g.Enable == nil {
		// return len(g.svr.GetServiceInfo()) > 0
		return len(g.svr.GetServiceInfo()) >= 0
	}
	return *g.Enable
}

func (g *Grpc) AddInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	g.interceptors = append(g.interceptors, interceptors...)
}
