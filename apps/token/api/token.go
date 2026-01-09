package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt/v5"
	authModel "github.com/kade-chen/google-billing-console/apps/common/model/auth"
	model "github.com/kade-chen/google-billing-console/apps/common/model/token"
	"github.com/kade-chen/google-billing-console/apps/token"
	tools "github.com/kade-chen/google-billing-console/tools/time"
	"github.com/kade-chen/google-billing-console/tools/trances"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/restful/response"
)

var LocAsiaShanghai, _ = time.LoadLocation("Asia/Shanghai")

func (h *tokenHandler) IssueToken(r *restful.Request, w *restful.Response) {
	req := token.NewIssueTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	// 补充用户的登录时的位置信息
	req.Location = token.NewNewLocationFromHttp(r.Request)
	tk, err := h.service.IssueToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}

	tk.SetCookie(w)
	response.Success(w, tk)
	// fmt.Println("success", tk)
	// 使用301永久重定向
}

func (h *tokenHandler) RevolkToken(r *restful.Request, w *restful.Response) {
	claims, ok := r.Attribute("claims").(*authModel.TokenAuthMiddleware)
	if !ok {
		response.Failed(w, exception.NewInternalServerError("unauthorized"))
		return
	}
	req := token.NewRevolkTokenRequest(claims.JwtToken, "")
	//3.delete cooke the token for mongodb
	ins, row, err := h.service.RevolkToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	ins.DeleteCookie(w)
	//5.log info
	// h.log.Info().Msgf("token revoke success, token: %s", ins.AccessToken)
	//6.return
	var result model.RevolkToken
	result.Action = "Revoke"
	result.AffectedRows = row
	result.Success = true
	response.Success(w, result)
}

func (h *tokenHandler) Validate_Token(r *restful.Request, w *restful.Response) {
	traceID := trances.NewTraceID()
	h.log.Info().Msg("validate token")
	req := token.NewValidateTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}
	//3.validate token
	ins, err := h.jwt.ValicateToken(traceID, req.AccessToken)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			var row model.ValidateToken
			row.EnableExpire = true
			row.JwtToken = req.AccessToken
			row.ExpiresAt = ins.ExpiresAt.Time.In(LocAsiaShanghai).Format("2006-01-02 15:04:05")
			response.Success(w, row)

		case errors.Is(err, jwt.ErrSignatureInvalid):
			// 伪造 token
			response.Failed(w, exception.NewUnauthorized("invalid signature"))

		default:
			response.Failed(w, exception.NewUnauthorized("invalid token"))
		}
		return
	}
	var row model.ValidateToken
	row.JwtToken = ins.JwtToken
	row.ExpiresAt = ins.ExpiresAt.Time.In(LocAsiaShanghai).Format("2006-01-02 15:04:05")

	response.Success(w, row)
}

func (h *tokenHandler) Refresh_Token(r *restful.Request, w *restful.Response) {
	trancesID := trances.NewTraceID()
	//1.获取cookis
	refresh_token, err := r.Request.Cookie("CCK")
	if err != nil {
		if err == http.ErrNoCookie {
			// 没有这个 cookie
			http.Error(w, "no CCK cookie", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//2.验证refresh_token是否过期
	tokenAuthMiddleware, err := h.jwt.ValicateToken(trancesID, refresh_token.Value)
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenExpired):
			// access token 过期 → refresh
			response.Failed(w, exception.NewUnauthorized("token expired"))
			// Redirect_Url(resp, err, t.log)

		case errors.Is(err, jwt.ErrSignatureInvalid):
			// 伪造 token
			response.Failed(w, exception.NewUnauthorized("invalid signature"))

		default:
			response.Failed(w, exception.NewUnauthorized("invalid token"))
		}
		return
	}
	endtime := tools.JwtTime(tokenAuthMiddleware.ExpiresAt)
	// fmt.Println(format.ToJSON(tokenAuthMiddleware))
	// fmt.Println("000000", endtime)
	// _ = tokenAuthMiddleware.ExpiresAt
	//3.获取access_token
	jwtToken, err := h.jwt.JwtRefreshAccessToken(tokenAuthMiddleware.Platform, tokenAuthMiddleware.Subject, endtime, tokenAuthMiddleware.Organizations)
	if err != nil {
		response.Failed(w, err)
		return
	}
	// 3.validate token/update access token
	// ins, err := h.service.UpdateToken(r.Request.Context(), &token.UpdateTokenRequest{
	// 	AccessToken:  jwtToken,
	// 	RefreshToken: refresh_token.Value,
	// })
	// if err != nil {
	// 	response.Failed(w, err)
	// 	// w.AddHeader("Location", "http://localhost:5173/login")
	// 	// w.WriteHeader(http.StatusFound) // 302
	// 	return
	// }
	response.Success(w, jwtToken)
}

func (h *tokenHandler) test(r *restful.Request, w *restful.Response) {
	_ = w.WriteJson("test", "application/json")
}
