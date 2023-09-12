//go:build manual

package integration

import (
	"database/sql"
	"encoding/json"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/integration/startup"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat"
	wechatmocks "gitee.com/geekbang/basic-go/webook/internal/service/oauth2/wechat/mocks"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// 这是一个只能手动运行的测试，为了摆脱 wechat 那个部分而引入的测试
func TestWechatCallback(t *testing.T) {
	const callbackUrl = "/oauth2/wechat/callback"
	db := startup.InitTestDB()
	testCases := []struct {
		name   string
		mock   func(ctrl *gomock.Controller) wechat.Service
		before func(t *testing.T)
		// 验证并且删除数据
		after      func(t *testing.T)
		wantCode   int
		wantResult web.Result
	}{
		{
			name: "注册新用户",
			mock: func(ctrl *gomock.Controller) wechat.Service {
				svc := wechatmocks.NewMockService(ctrl)
				svc.EXPECT().VerifyCode(gomock.Any(),
					gomock.Any()).
					Return(domain.WechatInfo{
						OpenId:  "123",
						UnionId: "1234",
					}, nil)
				return svc
			},
			before: func(t *testing.T) {
				// 什么也不需要做
			},
			after: func(t *testing.T) {
				// 验证数据库
				var u dao.User
				err := db.Find(&u, "wechat_open_id = ?", "123").Error
				assert.NoError(t, err)
				// 只需要验证 union id 就差不多了
				assert.Equal(t, "1234", u.WechatUnionId.String)
				db.Delete(&u, "wechat_open_id = ?", "123")
			},
			wantCode: 200,
			wantResult: web.Result{
				Msg: "登录成功",
			},
		},
		{
			name: "已有的用户",
			mock: func(ctrl *gomock.Controller) wechat.Service {
				svc := wechatmocks.NewMockService(ctrl)
				svc.EXPECT().VerifyCode(gomock.Any(), gomock.Any()).
					Return(domain.WechatInfo{
						OpenId:  "2345",
						UnionId: "23456",
					}, nil)
				return svc
			},
			before: func(t *testing.T) {
				// 插入数据，假装用户存在
				err := db.Create(&dao.User{
					WechatOpenId: sql.NullString{
						String: "2345",
						Valid:  true,
					},
					WechatUnionId: sql.NullString{
						String: "23456",
						Valid:  true,
					},
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var u dao.User
				err := db.Find(&u, "wechat_open_id = ?", "2345").Error
				assert.NoError(t, err)
				// 只需要验证 union id 就差不多了
				assert.Equal(t, "23456", u.WechatUnionId.String)
				db.Delete(&u, "wechat_open_id = ?", "123")
			},
			wantCode: 200,
			wantResult: web.Result{
				Msg: "登录成功",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			tc.before(t)

			userSvc := startup.InitUserSvc()
			jwtHdl := startup.InitJwtHdl()
			wechatSvc := tc.mock(ctrl)
			hdl := web.NewOAuth2WechatHandler(wechatSvc, userSvc, jwtHdl)
			server := gin.Default()
			hdl.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodGet,
				callbackUrl, nil)
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)

			code := recorder.Code
			// 反序列化为结果
			var result web.Result
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantCode, code)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}
