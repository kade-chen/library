package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/kade-chen/google-billing-console/apps/token"
	"github.com/kade-chen/library/exception"
	"github.com/kade-chen/library/http/restful/response"
)

func (h *tokenHandler) IssueToken(r *restful.Request, w *restful.Response) {
	req := token.NewIssueTokenRequest()
	if err := r.ReadEntity(req); err != nil {
		response.Failed(w, err)
		return
	}

	// 补充用户的登录时的位置信息
	// req.Location = token.NewNewLocationFromHttp(r.Request)
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
	req := token.NewRevolkTokenRequest("", "")
	//1.cookis
	if len(r.Request.Cookies()) != 0 {
		for _, v := range r.Request.Cookies() {
			req.ACCESS_TOKEN_NAME = v.Name
			req.AccessToken = v.Value
		}
		//2.header cookie
	} else if len(r.Request.URL.Query()) != 0 {

		for k, v := range r.Request.URL.Query() {
			req.ACCESS_TOKEN_NAME = k
			req.AccessToken = v[0]
		}
	} else {
		response.Failed(w, exception.NewAccessTokenIllegal("token is nil"))
		return
	}

	//3.delete cooke the token for mongodb
	ins, err := h.service.RevolkToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		return
	}
	//4.delete ui cookie
	ins.DeleteCookie(w)
	//5.log info
	h.log.Info().Msgf("token revoke success, token: %s", ins.AccessToken)
	//6.return
	response.Success(w, ins)
}

func (h *tokenHandler) Validate_Token(r *restful.Request, w *restful.Response) {
	h.log.Info().Msg("validate token")
	req := token.NewValidateTokenRequest("")
	//1.cookis
	if len(r.Request.Cookies()) != 0 {
		for _, v := range r.Request.Cookies() {
			req.ACCESS_TOKEN_NAME = v.Name
			req.AccessToken = v.Value
		}
		//2.header cookie
	} else if len(r.Request.URL.Query()) != 0 {
		for k, v := range r.Request.URL.Query() {
			req.ACCESS_TOKEN_NAME = k
			if req.ACCESS_TOKEN_NAME != "ACCESS_TOKEN_COOKIE_KEY" {
				response.Failed(w, exception.NewAccessTokenIllegal("token_name is error"))
				return
			}
			req.AccessToken = v[0]
		}
	} else {
		response.Failed(w, exception.NewAccessTokenIllegal("token is nil"))
		return
	}
	//3.validate token
	ins, err := h.service.ValicateToken(r.Request.Context(), req)
	if err != nil {
		response.Failed(w, err)
		// w.AddHeader("Location", "http://localhost:5173/login")
		// w.WriteHeader(http.StatusFound) // 302
		return
	}
	response.Success(w, ins)
}
