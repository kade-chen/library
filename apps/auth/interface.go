package auth

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/kade-chen/google-billing-console/apps/common/model/auth"
)

const (
	AppName = "TokenAuthMiddleware"
)

type Service interface {
	//接口认证
	Auth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain)
	//生成jwt token
	GeneratJwtAccessToken(platform int32, subject string, issueAt, ExpiredAt int64, organizations []string) (string, error)

	JwtRefreshAccessToken(platform int32, subject string, expiredAt int64, organizations []string) (accesstoken string, err error)

	//验证jwt token ,不做任何bq的处理
	ValicateToken(traceID, jwtToken string) (*auth.TokenAuthMiddleware, error)
}
