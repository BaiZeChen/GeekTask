package service

import "context"

type CodeService interface {
	Send(ctx context.Context,
		// 区别业务场景
		biz string, phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}
