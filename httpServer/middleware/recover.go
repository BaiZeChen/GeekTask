package middleware

import (
	"GeekTask/httpServer/route"
	"fmt"
)

type recoverBuildMiddleWare struct {
	statusCode int
}

func NewRecover(code int) *recoverBuildMiddleWare {
	return &recoverBuildMiddleWare{
		statusCode: code,
	}
}

func (r *recoverBuildMiddleWare) Build() route.Middleware {
	return func(next route.HandleFunc) route.HandleFunc {
		return func(ctx *route.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespStatusCode = r.statusCode
					ctx.RespData = []byte(fmt.Sprintf("系统错误，原因：%s", err))
				}
			}()
			next(ctx)
		}
	}
}
