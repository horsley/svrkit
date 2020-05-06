package svrkit

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	//ErrAuthInfoInvalid token 不合法
	ErrAuthInfoInvalid = errors.New("token invalid")
	//ErrAuthInfoExpired token 合法但已经过期
	ErrAuthInfoExpired = errors.New("token expired")
)

//CookieAuth 通过 cookie 存放登录态+签名校验，通过添加前置登录态检查，把没有登录态的请求重定向到登录页面，附带过期校验
type CookieAuth struct {
	//LoginPage 没有登录态的跳转页面
	LoginPage string

	//CookieName 存放cookie的 key
	CookieName string

	//TTL 登录态有效期，默认永久有效
	TTL time.Duration

	//Cookie 存储实现，可设置签名盐值
	Store CookieStore

	//自定义登录态校验方法，返回 false 为校验失败重定向到跳转页面
	VerifyFunc func(info string, authExpired bool) bool
}

//GuardFunc 为http.HandlerFunc 添加前置登录态检查
func (c *CookieAuth) GuardFunc(h http.HandlerFunc) http.HandlerFunc {
	return c.GuardHandler(http.Handler(h)).ServeHTTP
}

//Guard 为 svrkit.HTTPHandlerFunc 添加前置登录态检查
func (c *CookieAuth) Guard(fn HTTPHandlerFunc) http.HandlerFunc {
	return c.GuardHandler(http.Handler(HTTPFunc(fn))).ServeHTTP
}

//GuardHandler 为http.Handler 添加前置登录态检查
func (c *CookieAuth) GuardHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		info, err := c.GetAuthInfo(r)
		if c.VerifyFunc != nil && !c.VerifyFunc(info, err == ErrAuthInfoExpired) {
			http.Redirect(rw, r, c.LoginPage, http.StatusFound)
			return
		}

		h.ServeHTTP(rw, r)
	})
}

//GuardRouter 为 svrkit.Router 整体提供保护
func (c *CookieAuth) GuardRouter(r *Router) *Router {
	wrap := NewRouter()

	wrap.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		info, err := c.GetAuthInfo(req)
		if c.VerifyFunc != nil && !c.VerifyFunc(info, err == ErrAuthInfoExpired) {
			http.Redirect(rw, req, c.LoginPage, http.StatusFound)
			return
		}

		r.ServeHTTP(rw, req)
	})

	return wrap
}

// SetAuthInfo 写入登录态信息
func (c *CookieAuth) SetAuthInfo(rw http.ResponseWriter, info string) {
	if c.TTL.Nanoseconds() == 0 {
		c.Store.SetWithTTL(rw, c.CookieName, c.prependTimeInfo(info), 365*10*24*60*60*time.Second) //10 years for forever
	} else {
		c.Store.SetWithTTL(rw, c.CookieName, c.prependTimeInfo(info), c.TTL*2)
	}
}

func (c *CookieAuth) prependTimeInfo(info string) string {
	ts := time.Now().Unix()
	return fmt.Sprintf("%d|%s", ts, info)
}

//GetAuthInfo 从cookie中获取登录态信息，注意当错误为过期的时候，信息和 error 同时有值，使用者可以决定后续策略（包括自动续期、判为过期）
func (c *CookieAuth) GetAuthInfo(r *http.Request) (string, error) {
	val := c.Store.Get(r, c.CookieName)
	if val != "" {
		parts := strings.SplitN(val, "|", 2)
		if len(parts) == 2 {
			if c.TTL.Nanoseconds() != 0 {
				ts, _ := strconv.Atoi(parts[1])
				if time.Now().Unix()-int64(ts) > int64(c.TTL) {
					return parts[1], ErrAuthInfoExpired
				}
			}

			return parts[1], nil
		}
	}

	return "", ErrAuthInfoInvalid
}
