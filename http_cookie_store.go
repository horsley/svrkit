package svrkit

import (
	"errors"
	"net/http"
)

//CookieStore 在 cookie 中存放会话信息并附带签名，可防伪造
type CookieStore struct {
	//SignSalt 签名的盐
	SignSalt string

	//Path cookie 关联路径，rest形式url需要注意cookie设置路径，否则影响读取
	Path string
}

const sha1HashLen = 40

//ErrSignVerifyFail 校验签名失败
var ErrSignVerifyFail = errors.New("sign not match")

//Set 写入会话信息，注意 cookie 使用须在内容 body 输出之前
func (c *CookieStore) Set(rw http.ResponseWriter, k, v string) {
	http.SetCookie(rw, &http.Cookie{
		Name:  k,
		Value: c.signValue(v),
		Path:  c.Path,
	})
}

//SetCookie 写入cookie信息，方便使用方调整 cookie 条目属性，本方法只会修改 cookie 值附上签名
func (c *CookieStore) SetCookie(rw http.ResponseWriter, cookie *http.Cookie) {
	if cookie != nil {
		cookieCopy := *cookie

		cookieCopy.Value = c.signValue(cookieCopy.Value)
		http.SetCookie(rw, &cookieCopy)
	}
}

//Get 读取会话信息，会校验签名，失败情况返回空串
func (c *CookieStore) Get(r *http.Request, k string) string {
	if cookie, err := r.Cookie(k); err == nil {
		if val, err := c.checkSign(cookie.Value); err == nil {
			return val
		}
	}
	return ""
}

func (c *CookieStore) signValue(v string) string {
	sign := SHA1Hash(c.SignSalt + SHA1Hash(v))
	return sign + v
}

func (c *CookieStore) checkSign(origin string) (string, error) {
	if len(origin) >= sha1HashLen {
		val := origin[sha1HashLen:]

		if c.signValue(val) == origin {
			return val, nil
		}
	}
	return "", ErrSignVerifyFail
}
