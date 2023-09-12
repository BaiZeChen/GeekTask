package web

import (
	"GeekTask/fifthWeek/internal/service"
	"GeekTask/fifthWeek/internal/service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_LoginSMS(t *testing.T) {
	const signupUrl = "/users/login_sms"
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		wantCode   int
		wantBody   string
	}{
		{
			name: "codeSvc系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				codeSvc := mocks.NewMockCodeService(ctrl)
				codeSvc.EXPECT().
					Verify(gomock.Any(), gomock.Eq(biz), gomock.Eq("1101121131141"), gomock.Eq("123456")).
					Return(false, errors.New("随便一个报错"))
				return nil, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				body := bytes.NewBuffer([]byte(`{"phone":"1101121131141","code":"123456"}`))
				req, err := http.NewRequest(http.MethodPost, signupUrl, body)
				req.Header.Set("Content-Type", "application/json")
				if err != nil {
					t.Fatal(err)
				}
				return req
			},
			wantCode: 5,
			wantBody: "系统错误",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctl := gomock.NewController(t)
			defer ctl.Finish()

			uSvc, codeSvc := tc.mock(ctl)
			hdl := NewUserHandler(uSvc, codeSvc)

			ginSvc := gin.Default()
			hdl.RegisterRoutes(ginSvc)
			// 准备请求
			req := tc.reqBuilder(t)
			// 准备记录响应
			recorder := httptest.NewRecorder()
			// 执行
			ginSvc.ServeHTTP(recorder, req)
			result := &Result{}
			json.Unmarshal(recorder.Body.Bytes(), result)
			// 断言
			assert.Equal(t, tc.wantCode, result.Code)
			assert.Equal(t, tc.wantBody, result.Msg)
		})
	}
}
