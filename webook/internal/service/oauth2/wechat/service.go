package wechat

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2"
	"golang.org/x/net/context"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

type Service struct {
	appId string
}

func NewService(appId string) oauth2.Service {
	return &Service{
		appId: appId,
	}
}

func (s *Service) AuthURL(ctx context.Context) (string, error) {
	//stats := uuid
	//return
	panic("implement me")
}
