package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

var _ handler = (*OAuth2WechatHandler)(nil)

type OAuth2WechatHandler struct {
	// 这边也可以直接定义成 wechat.Service
	// 但是为了保持使用 mock 来测试，这里还是用了接口
	wechatSvc wechat.Service
	userSvc   service.UserService
	jwtHandler
}

func NewOAuth2WechatHandler(service wechat.Service,
	userSvc service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatSvc: service,
		userSvc:   userSvc,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.OAuth2URL)
	// 这边用 Any 万无一失
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	state := ctx.Query("state")
	code := ctx.Query("code")
	info, err := h.wechatSvc.VerifyCode(ctx, code, state)
	if err != nil {
		// 实际上这个错误，也有可能是 code 不对
		// 但是给前端的信息没有太大的必要区分究竟是代码不对还是系统本身有问题
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	// 这里就是登录成功
	// 所以你需要设置 JWT
	u, err := h.userSvc.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = h.setJWTToken(ctx, u.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

func (h *OAuth2WechatHandler) OAuth2URL(ctx *gin.Context) {
	url, err := h.wechatSvc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误，请稍后再试",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
	return
}
