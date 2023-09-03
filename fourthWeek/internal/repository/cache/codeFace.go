package cache

import (
	"context"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
)

type CodeCache interface {
	Set(ctx context.Context, lockKey, bizKey, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

func NewCodeCacheFace(redisClient redis.Cmdable, localClient *freecache.Cache) CodeCache {
	return &CodeCacheFace{
		redisCache: redisClient,
		localCache: localClient,
	}
}

type CodeCacheFace struct {
	redisCache redis.Cmdable
	localCache *freecache.Cache
}

func (c *CodeCacheFace) Set(ctx context.Context, lockKey, bizKey, phone, code string) error {
	//TODO implement me
	panic("implement me")
}

func (c *CodeCacheFace) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
