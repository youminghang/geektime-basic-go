package startup

import "gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService() wechat.Service {
	return wechat.NewService("", "")
}
