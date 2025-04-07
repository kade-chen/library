package exception

// -----------------------------400 code start --------------------------------------//
// NewBadRequest  400
// a 的作用 "Resource not found: %s", "exampleID"
func NewBadRequest(format string, a ...interface{}) *APIException {
	return NewAPIException("", BadRequest, codeReason(BadRequest), format, a...)
}

// Unauthorized 401
func NewUnauthorized(format string, a ...interface{}) *APIException {
	return NewAPIException("", Unauthorized, codeReason(Unauthorized), format, a...)
}

// Forbidden 403
func NewForbidden(format string, a ...interface{}) *APIException {
	return NewAPIException("", Forbidden, codeReason(Forbidden), format, a...)
}

// NewNotFound  404
// a 的作用 "Resource not found: %s", "exampleID"
func NewNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException("", NotFound, codeReason(NotFound), format, a...)
}

// Conflict 409 资源冲突, 已经存在
func NewConflict(format string, a ...interface{}) *APIException {
	return NewAPIException("", Conflict, codeReason(Conflict), format, a...)
}

// NewInternalServerError 500
func NewInternalServerError(format string, a ...interface{}) *APIException {
	return NewAPIException("", InternalServerError, codeReason(InternalServerError), format, a...)
}

//-----------------------------400-500 code end --------------------------------------//

// -----------------------------50000-9999 code start --------------------------------------//
// OtherPlaceLoggedIn 50010 异地登录
func NewOtherPlaceLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException("", OtherPlaceLoggedIn, codeReason(OtherPlaceLoggedIn), format, a...)
}

// OtherIPLoggedIn 50011 异常IP登录
func NewOtherIPLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException("", OtherIPLoggedIn, codeReason(OtherIPLoggedIn), format, a...)
}

// OtherClientsLoggedIn 50012 用户已经通过其他端登录
func NewOtherClientsLoggedIn(format string, a ...interface{}) *APIException {
	return NewAPIException("", OtherClientsLoggedIn, codeReason(OtherClientsLoggedIn), format, a...)
}

// SessionTerminated 50013 会话中断
func NewSessionTerminated(format string, a ...interface{}) *APIException {
	return NewAPIException("", SessionTerminated, codeReason(SessionTerminated), format, a...)
}

// AccessTokenExpired 50014 token过期
func NewAccessTokenExpired(format string, a ...interface{}) *APIException {
	return NewAPIException("", AccessTokenExpired, codeReason(AccessTokenExpired), format, a...)
}

// RefreshTokenExpired 50015 token过期
func NewRefreshTokenExpired(format string, a ...interface{}) *APIException {
	return NewAPIException("", RefreshTokenExpired, codeReason(RefreshTokenExpired), format, a...)
}

// RefreshTokenExpired 50016 访问token不合法
func NewAccessTokenIllegal(format string, a ...interface{}) *APIException {
	return NewAPIException("", AccessTokenIllegal, codeReason(AccessTokenIllegal), format, a...)
}

// RefreshTokenIllegal 50017 刷新token不合法
func NewRefreshTokenIllegal(format string, a ...interface{}) *APIException {
	return NewAPIException("", RefreshTokenIllegal, codeReason(RefreshTokenIllegal), format, a...)
}

// VerifyCodeRequired 50018  需要验证码
func NewVerifyCodeRequired(format string, a ...interface{}) *APIException {
	return NewAPIException("", VerifyCodeRequired, codeReason(VerifyCodeRequired), format, a...)
}

// VerifyCodeRequired 50019  用户密码过期
func NewPasswordExired(format string, a ...interface{}) *APIException {
	return NewAPIException("", PasswordExired, codeReason(PasswordExired), format, a...)
}

// NewWebCookisNotFoundnRequest 50020 cookis is not found for web 浏览器未发现cookis
// a 的作用 "Resource not found: %s", "exampleID"
func NewWebCookisNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException("", WebCookisNotFound, codeReason(WebCookisNotFound), format, a...)
}

// ApiUrlNotFound 50040 api is not found for web
func NewApiUrlNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException("", ApiUrlNotFound, codeReason(ApiUrlNotFound), format, a...)
}

// PermissionAuthenticationFailed 50041 permission authentication failed
func NewPermissionAuthenticationFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", PermissionAuthenticationFailed, codeReason(PermissionAuthenticationFailed), format, a...)
}

// IocRegisterFailed 50050 ioc register failed
func NewIocRegisterFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", IocRegisterFailed, codeReason(IocRegisterFailed), format, a...)
}

// IocImplRegisterFailed 50051 ioc impl register failed
func NewIocImplRegisterFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", IocImplRegisterFailed, codeReason(IocImplRegisterFailed), format, a...)
}

// IocApiRegisterFailed 50052 ioc get bean failed
func NewIocApiRegisterFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", IocApiRegisterFailed, codeReason(IocApiRegisterFailed), format, a...)
}

// IocGetFailed 50053 ioc get bean failed
func NewIocGetFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", IocGetFailed, codeReason(IocGetFailed), format, a...)
}

// ProviderNotFound 50054 provider not found
func NewProviderNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderNotFound, codeReason(ProviderNotFound), format, a...)
}

// ProviderRegistryNil 50055 provider registry nil
func NewProviderRegistryNil(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderRegistryNil, codeReason(ProviderRegistryNil), format, a...)
}

// ProviderRegistryFailed 50056 provider registry failed
func NewProviderRegistryFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderRegistryFailed, codeReason(ProviderRegistryFailed), format, a...)
}

// ProviderTokenRegistryNil 50057 token provider registry nil
func NewProviderTokenRegistryNil(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderTokenRegistryNil, codeReason(ProviderTokenRegistryNil), format, a...)
}

// ProviderTokenRegistryFailed 50058 token provider registry failed
func NewProviderTokenRegistryFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderTokenRegistryFailed, codeReason(ProviderTokenRegistryFailed), format, a...)
}

// ProviderVertexRegistryNil 50059 vertex provider registry nil
func NewProviderTokenNotFound(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderVertexRegistryNil, codeReason(ProviderVertexRegistryNil), format, a...)
}

// ProviderVertexRegistryFailed 50060 vertex provider registry failed
func NewProviderVertexRegistryFailed(format string, a ...interface{}) *APIException {
	return NewAPIException("", ProviderVertexRegistryFailed, codeReason(ProviderVertexRegistryFailed), format, a...)
}

// UnKnownException 99999 未知异常
func NewUnKnownException(format string, a ...interface{}) *APIException {
	return NewAPIException("", UnKnownException, codeReason(UnKnownException), format, a...)
}

//-----------------------------50000-9999 code end --------------------------------------//
