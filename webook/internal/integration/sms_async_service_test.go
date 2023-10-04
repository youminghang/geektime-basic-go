package integration

import (
	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/service/sms"
	smsmocks "gitee.com/geekbang/basic-go/webook/internal/service/sms/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"testing"
)

type AsyncSMSTestSuite struct {
	suite.Suite
	db *gorm.DB
}

func (s *AsyncSMSTestSuite) SetupSuite() {
	s.db = startup.InitTestDB()
}

func (s *AsyncSMSTestSuite) TestSend() {
	t := s.T()
	testCases := []struct {
		name string

		// 虽然是集成测试，但是我们也不想真的发短信，所以用 mock
		mock func(ctrl *gomock.Controller) sms.Service

		tplId   string
		args    []string
		numbers []string

		wantErr error
	}{
		{
			name: "异步",
			mock: func(ctrl *gomock.Controller) sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				return svc
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := startup.InitAsyncSmsService(tc.mock(ctrl))
			err := svc.Send(context.Background(), tc.tplId, tc.args, tc.numbers...)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func (s *AsyncSMSTestSuite) TestSend() {

}

func TestAsyncSmsService(t *testing.T) {
	suite.Run(t, &AsyncSMSTestSuite{})
}
