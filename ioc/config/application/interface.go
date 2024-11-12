package application

import "github.com/kade-chen/library/ioc"

const (
	AppName = "app"
)

func Get() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
