package ioc

import (
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/localsms"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms/tencent"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"os"
)

func InitSmsService() sms.Service {
	//return initSmsTencentService()
	return InitSmsMemoryService()
}

func initSmsTencentService() sms.Service {
	// 在这里你也可以考虑从配置文件里面读取
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("没有找到环境变量 SMS_SECRET_ID ")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("没有找到环境变量 SMS_SECRET_KEY")
	}

	c, err := tencentSMS.NewClient(common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile())
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "妙影科技")
}

// InitSmsMemoryService 使用基于内存，输出到控制台的实现
func InitSmsMemoryService() sms.Service {
	return localsms.NewService()
}
