package log

import (
	"github.com/kade-chen/library/ioc"
	"github.com/rs/zerolog"
)

const (
	AppName        = "log"
	SUB_LOGGER_KEY = "component"
)

// 没有key:value 直接调用直走init
func L() *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Config).root
}

func Sub(name string) *zerolog.Logger {
	//return ioc.Config().Get(AppName).(*Config).root,这个key：value可以在init里面做也就是L里面的
	return ioc.Config().Get(AppName).(*Config).Logger(name)
}

func T(name string) *TraceLogger {
	//"level":"debug"
	return ioc.Config().Get(AppName).(*Config).TraceLogger(name)
}
