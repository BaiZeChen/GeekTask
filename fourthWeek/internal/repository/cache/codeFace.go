package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/redis/go-redis/v9"
	"strconv"
	"sync"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode             string
	ErrCodeSendTooMany        = errors.New("发送验证码太频繁")
	ErrUnknownForCode         = errors.New("发送验证码遇到未知错误")
	ErrCodeVerifyTooManyTimes = errors.New("验证次数太多")
)

type CodeCache interface {
	Set(ctx context.Context, bizKey, phone, code string) error
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

func (c *CodeCacheFace) Set(ctx context.Context, bizKey, phone, code string) error {
	res, err := c.redisCache.Eval(ctx, luaSetCode, []string{c.key(bizKey, phone)}, code).Int()
	if err != nil {
		// 认为redis不能用了
		return c.localSet(bizKey, phone, code)
	}
	switch res {
	case 0:
		go func() {
			// 此时，本地缓存也存一份并忽略错误，防止get不能用
			_ = c.localSet(bizKey, phone, code)
		}()
		return nil
	case -1:
		//	最近发过
		return ErrCodeSendTooMany
	default:
		// 系统错误，比如说 -2，是 key 冲突
		// 其它响应码，不知道是啥鬼东西
		// TODO 按照道理，这里要考虑记录日志，但是我们暂时还没有日志模块，所以暂时不管
		return ErrUnknownForCode
	}
}

func (c *CodeCacheFace) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	res, err := c.redisCache.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		// 认为redis不能用了
		return c.localVerify(biz, phone, inputCode)
	}

	go func() {
		// 不管那个分支，需要同步一下本地缓存，可以忽略错误
		_, _ = c.localVerify(biz, phone, inputCode)
	}()
	switch res {
	case 0:
		return true, nil
	case -1:
		//	验证次数耗尽，一般都是意味着有人在捣乱
		return false, ErrCodeVerifyTooManyTimes
	default:
		// 验证码不对
		return false, nil
	}
}

func (c *CodeCacheFace) localSet(bizKey, phone, code string) error {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	// 允许与redis过期时间有误差
	ttl, err := c.localCache.TTL([]byte(c.key(bizKey, phone)))
	if err != nil || ttl < 540 {
		// 代表可以重发
		err = c.localCache.Set([]byte(c.key(bizKey, phone)), []byte(code), 600)
		if err != nil {
			return ErrUnknownForCode
		}
		err = c.localCache.Set([]byte(c.key(bizKey, phone)+":cnt"), []byte("3"), 600)
		if err != nil {
			return ErrUnknownForCode
		}
		return nil
	}
	return ErrCodeSendTooMany
}

func (c *CodeCacheFace) localVerify(bizKey, phone, inputCode string) (bool, error) {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	code, err := c.localCache.Get([]byte(c.key(bizKey, phone)))
	if err != nil {
		return false, err
	}
	cntStr, ttl, err := c.localCache.GetWithExpiration([]byte(c.key(bizKey, phone) + ":cnt"))
	if err != nil {
		return false, err
	}
	cnt, _ := strconv.Atoi(string(cntStr))
	if cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	if string(code) == inputCode {
		_ = c.localCache.Set([]byte(c.key(bizKey, phone)+":cnt"), []byte("-1"), int(ttl))
		return true, nil
	}

	// 用户输错了
	cnt--
	cntStr = []byte(strconv.Itoa(cnt))
	_ = c.localCache.Set([]byte(c.key(bizKey, phone)+":cnt"), cntStr, int(ttl))
	return false, nil
}

func (c *CodeCacheFace) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
