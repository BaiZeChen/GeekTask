package middleware

import "GeekTask/httpServer/route"

type flashRespBuildMiddleWare struct{}

func NewFlashResp() *flashRespBuildMiddleWare {
	return &flashRespBuildMiddleWare{}
}

func (f *flashRespBuildMiddleWare) Build() route.Middleware {
	return func(next route.HandleFunc) route.HandleFunc {
		return func(ctx *route.Context) {
			next(ctx)
			if ctx.RespStatusCode > 0 {
				ctx.Resp.WriteHeader(ctx.RespStatusCode)
			}
			ctx.Resp.Write(ctx.RespData)
		}
	}
}
