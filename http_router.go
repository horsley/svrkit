package svrkit

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

//Router 是 http.ServeMux 增强版
type Router struct {
	*http.ServeMux

	//BeforeHandler router 的前置过滤器 可以通过返回 true 跳过后续处理（包括 AfterHandler）拦截请求
	BeforeHandler func(*ResponseWriter, *Request) (skip bool)

	//AfterHandler router 的后置处理器，可以从 httptest.ResponseRecorder中取得 handle 过程输出的 header 和 responseBody
	//如返回 skip = true，则不做后续处理，否则会把 recorder 录得的内容在 ResponseWriter 中重放，相当于 AfterHandler 透明
	AfterHandler func(*ResponseWriter, *Request, *httptest.ResponseRecorder) (skip bool)
}

//requestUserInfo 这里暂存 beforeHandler 中可能写入的 userInfo, 否则HTTPFuncHandler转成http 包标准 handler 之后，丢失userInfo 信息
var requestUserInfo = make(map[string]interface{})

//NewRouter 创建 router 的方法
func NewRouter() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
	}
}

//HandleFuncEx 本包 handler 的注册方法
func (rt *Router) HandleFuncEx(pattern string, handler HTTPHandlerFunc) {
	rt.HandleFunc(pattern, HTTPFunc(handler))
}

//SetSubRouter 设置二级的路由
func (rt *Router) SetSubRouter(pattern string, sub *Router) {
	rt.Handle(pattern, sub)
}

func (rt *Router) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if rt.BeforeHandler != nil {
		req := &Request{Request: r}
		skip := rt.BeforeHandler(&ResponseWriter{rw}, req)
		if skip {
			return
		}
		if req.UserInfo != nil {
			requestUserInfo[fmt.Sprintf("%p", r)] = req.UserInfo
		}
	}

	if rt.AfterHandler != nil {
		rec := httptest.NewRecorder()
		rt.ServeMux.ServeHTTP(rec, r)
		skip := rt.AfterHandler(&ResponseWriter{rw}, &Request{Request: r}, rec)
		if skip {
			return
		}

		for k := range rec.HeaderMap {
			rw.Header().Set(k, rec.HeaderMap.Get(k))
		}

		if rec.Code != 0 {
			rw.WriteHeader(rec.Code)
		}

		rw.Write(rec.Body.Bytes())
		return
	}

	rt.ServeMux.ServeHTTP(rw, r)
}
