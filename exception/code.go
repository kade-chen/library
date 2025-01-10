package exception

import "net/http"

const (
	// OtherPlaceLoggedIn 异地登录
	OtherPlaceLoggedIn = 50010
	// OtherIPLoggedIn 异常IP登录
	OtherIPLoggedIn = 50011
	// OtherClientsLoggedIn 用户已经通过其他端登录
	OtherClientsLoggedIn = 50012
	// SessionTerminated 会话中断
	SessionTerminated = 50013
	// AccessTokenExpired token过期
	AccessTokenExpired = 50014
	// RefreshTokenExpired token过期
	RefreshTokenExpired = 50015
	// AccessTokenIllegal 访问token不合法
	AccessTokenIllegal = 50016
	// RefreshTokenIllegal 刷新token不合法
	RefreshTokenIllegal = 50017
	// VerifyCodeRequired 需要验证码
	VerifyCodeRequired = 50018
	// PasswordExired 用户密码过期
	PasswordExired = 50019
	// cookis is not found for web  浏览器未发现cookis
	WebCookisNotFound = 50020
	//api is not found for web
	ApiUrlNotFound = 50040

	//permission authentication failed
	PermissionAuthenticationFailed = 50041

	//IocRegisterFailed
	IocRegisterFailed = 50050
	//IocImplRegisterFailed
	IocImplRegisterFailed = 50051
	//IocApiRegisterFailed
	IocApiRegisterFailed = 50052
	//IocGetFailed
	IocGetFailed = 50053

	// BadRequest 请求不合法
	BadRequest = http.StatusBadRequest
	// Unauthorized 未认证
	Unauthorized = http.StatusUnauthorized
	// Forbidden 无权限
	Forbidden = http.StatusForbidden
	// NotFound 接口未找到
	NotFound = http.StatusNotFound
	// Conflict 资源冲突, 已经存在
	Conflict = http.StatusConflict
	// InternalServerError 服务端内部错误
	InternalServerError = http.StatusInternalServerError

	// UnKnownException 未知异常
	UnKnownException = 99999
)

var (
	reasonMap = map[int]string{
		Unauthorized:                   "认证失败",
		NotFound:                       "资源未找到",
		Conflict:                       "资源已经存在",
		BadRequest:                     "请求不合法",
		InternalServerError:            "系统内部错误",
		Forbidden:                      "访问未授权",
		UnKnownException:               "未知异常",
		AccessTokenIllegal:             "访问令牌不合法",
		RefreshTokenIllegal:            "刷新令牌不合法",
		OtherPlaceLoggedIn:             "异地登录",
		OtherIPLoggedIn:                "异常IP登录",
		OtherClientsLoggedIn:           "用户已经通过其他端登录",
		SessionTerminated:              "会话结束",
		AccessTokenExpired:             "访问过期, 请刷新",
		RefreshTokenExpired:            "刷新过期, 请登录",
		VerifyCodeRequired:             "异常操作, 需要验证码进行二次确认",
		PasswordExired:                 "密码过期, 请找回密码或者联系管理员重置",
		WebCookisNotFound:              "浏览器未发现cookis所带的access_tokens",
		ApiUrlNotFound:                 "api 不合法",
		PermissionAuthenticationFailed: "权限认证失败",
		IocRegisterFailed:              "ioc 注册失败",
		IocImplRegisterFailed:          "ioc 实现类注册失败",
		IocApiRegisterFailed:           "ioc api 注册失败",
		IocGetFailed:                   "ioc 获取失败",
	}
)

// 返回code码对应的信息
func codeReason(code int) string {
	v, ok := reasonMap[code]
	if !ok {
		v = reasonMap[UnKnownException]
	}

	return v
}
