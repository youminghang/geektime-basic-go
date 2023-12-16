package startup

import (
	intrv1 "gitee.com/geekbang/basic-go/webook/api/proto/gen/intr/v1"
	"gitee.com/geekbang/basic-go/webook/bff/client"
	"gitee.com/geekbang/basic-go/webook/interactive/service"
)

func InitInteractiveClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	return client.NewInteractiveServiceAdapter(svc)
}
