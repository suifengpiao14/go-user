package open

import (
	"github.com/chanxuehong/util/random"
	"github.com/chanxuehong/wechat/open/oauth2"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/aiyi/go-user/config"
	"github.com/aiyi/go-user/frontend/errors"
	"github.com/aiyi/go-user/frontend/session"
	"github.com/aiyi/go-user/frontend/token"
)

// 获取微信登录页面的 url
//  redirect_uri 需要提供, 相对路径
func LoginURLHandler(ctx *gin.Context) {
	// NOTE: 在此之前的中间件获取了 token 和 session
	queryValues := ctx.Request.URL.Query()
	redirectURI := queryValues.Get("redirect_uri")
	if redirectURI == "" {
		ctx.JSON(200, errors.ErrBadRequest)
		return
	}
	redirectURI = config.ConfigData.WebServer.BaseURL + redirectURI

	tk := ctx.MustGet("token").(*token.Token)
	ss := ctx.MustGet("session").(*session.Session)

	ss.OAuth2State = string(random.NewRandomEx())
	if err := session.Set(tk.SessionId, ss); err != nil {
		glog.Errorln(err)
		ctx.JSON(200, errors.ErrInternalServerError)
		return
	}

	loginURL := oauth2.AuthCodeURL(config.ConfigData.Weixin.Open.AppId, redirectURI, "snsapi_login", ss.OAuth2State, nil)

	resp := struct {
		*errors.Error
		URL string `json:"url"`
	}{
		Error: errors.ErrOK,
		URL:   loginURL,
	}
	ctx.JSON(200, &resp)
	return
}