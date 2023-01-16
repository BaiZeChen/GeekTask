package route

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req    *http.Request
	Resp   http.ResponseWriter
	Params map[string]string
	// 缓存的数据
	cacheQueryValues url.Values
}

func (c *Context) JsonBind(val interface{}) error {
	if c.Req.Body == nil {
		return errors.New("web body数据为空")
	}

	err := json.NewDecoder(c.Req.Body).Decode(val)
	if err != nil {
		return err
	}

	return nil

}

func (c *Context) FormValue(key string) *StringValue {
	if err := c.Req.ParseForm(); err != nil {
		return &StringValue{err: err}
	}
	return &StringValue{val: c.Req.FormValue(key)}
}

func (c *Context) QueryValue(key string) *StringValue {
	// 这里与c.Req.ParseForm()不一样
	// c.Req.ParseForm()重复调用时不会重新获取
	// c.Req.URL.Query()重复调用时每次都会重新获取
	// 所以这个咱们都设置一下缓存
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return &StringValue{err: errors.New("没有找到对应的值")}
	}
	return &StringValue{val: vals[0]}
}

func (c *Context) PathValue(key string) *StringValue {
	val, ok := c.Params[key]
	if !ok {
		return &StringValue{err: errors.New("没有找到对应的参数值")}
	}
	return &StringValue{val: val}
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}

func (c *Context) JsonResp(code int, val interface{}) error {
	res, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.Resp.WriteHeader(code)
	_, err = c.Resp.Write(res)
	if err != nil {
		return err
	}
	return nil
}

type StringValue struct {
	val string
	err error
}

func (s *StringValue) String() (string, error) {
	return s.val, s.err
}

func (s *StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
