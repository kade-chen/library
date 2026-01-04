package impl

import (
	"github.com/kade-chen/google-billing-console/apps/auth"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/log"
	"github.com/rs/zerolog"
)

var _ auth.Service = (*service)(nil)

func init() {
	ioc.Controller().Registry(&service{})
}

// 用于鉴权的中间件
// 用于Token鉴权的中间件
type service struct {
	tk   token.Service
	auth authModel.TokenAuthMiddleware
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (s *service) Init() error {
	s.log = log.Sub(s.Name())
	// tk: ioc.Default().Get(tokenimpl.AppName).(*tokenimpl.TokenServiceImpl),
	s.tk = ioc.Controller().Get(token.AppName).(token.Service)
	// user: ioc.Controller().Get(user.AppName).(user.Service),
	return nil
}

func (service) Name() string {
	return auth.AppName
}

func (*service) Priority() int {
	return -1
}
