package grpc

import "github.com/kade-chen/library/ioc"

const (
	AppName = "grpc"
)

func Get() *Grpc {
	return ioc.Config().Get(AppName).(*Grpc)
}
