package svrkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/alecthomas/template"
)

var (
	//CookieTokenPrivateKey CookieToken hash用的盐值
	CookieTokenPrivateKey = "3f47d8c41e5576r09h1310H)(*"
)

//HTTPSvr http 服务器
type HTTPSvr struct {
	http.Server
}

//ResponseWriter ResponseWriter 封装一些便捷操作
type ResponseWriter struct {
	http.ResponseWriter
}

//WriteString 输出字符串
func (rw *ResponseWriter) WriteString(msg string) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Write([]byte(msg))
}

//WriteJSON 输出 json
func (rw *ResponseWriter) WriteJSON(data interface{}) {
	b, _ := json.Marshal(data)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(b)
}

//RenderHTML 渲染网页 返回渲染结果
func (rw *ResponseWriter) RenderHTML(file string, data interface{}) []byte {
	t, err := template.ParseFiles(file)
	if err != nil {
		return []byte("tmpl parse error:" + err.Error())
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	if err != nil {
		return []byte("tmpl exec error:" + err.Error())
	}

	return buf.Bytes()
}

//SetCookieToken 种植 CookieToken
func (rw *ResponseWriter) SetCookieToken(key, val string) {
	http.SetCookie(rw, &http.Cookie{
		Name:  key,
		Value: fmt.Sprintf("%s|%s", val, SHA1Hash(CookieTokenPrivateKey+SHA1Hash(val))),
		Path:  "/",
	})
}

//Request HTTP请求，封装一些便捷操作
type Request struct {
	*http.Request
	reqBody []byte
}

//IsPost 是否 post 请求
func (r *Request) IsPost() bool {
	return r.Method == "POST"
}

//ClientIP 获取 ip 地址
func (r *Request) ClientIP() string {
	var resultIP string
	realIP := r.Header.Get("X-Real-Ip")
	fwFor := r.Header.Get("X-Forwarded-For")

	if realIP == "" && fwFor == "" { //无代理
		spIdx := strings.LastIndex(r.RemoteAddr, ":") //从尾部找冒号，因为IPv6在地址中部有冒号
		if spIdx == -1 {
			return r.RemoteAddr
		}
		return r.RemoteAddr[:spIdx]
	} else if fwFor != "" { //优先看forward for
		ipLst := strings.Split(fwFor, ",")
		resultIP = strings.TrimSpace(ipLst[0])
	} else {
		resultIP = realIP
	}
	//由于fw头可能是伪造的，这里进行一次校验
	if net.ParseIP(resultIP) != nil {
		return resultIP
	}
	return ""
}

//ReadJSON 解析请求里面的json
func (r *Request) ReadJSON(data interface{}) error {
	bin, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if len(bin) == 0 {
		return fmt.Errorf("empty req body")
	}

	return json.Unmarshal(bin, data)
}

//ReadRequestBody 读取请求body，这个方法不受 body 只能读一次的限制，我们会把内容存起来
func (r *Request) ReadRequestBody() []byte {
	if r.reqBody == nil {
		bin, err := ioutil.ReadAll(r.Body)
		if err == nil {
			r.reqBody = bin
		}
	}
	return r.reqBody
}

//GetCookieToken 获取存放在 cookie 的 kv 值，如果没有或者检验不合法返回空串
func (r *Request) GetCookieToken(key string) string {
	if c, err := r.Cookie(key); err == nil {
		parts := strings.Split(c.Value, "|")
		if len(parts) != 2 {
			return ""
		}

		val := parts[0]
		hash := parts[1]

		if SHA1Hash(CookieTokenPrivateKey+SHA1Hash(val)) != hash {
			return ""
		}

		return val
	}
	return ""
}

//HTTPFunc http 包装
func HTTPFunc(handler func(*ResponseWriter, *Request)) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		rspW := ResponseWriter{rw}
		req := Request{Request: r}
		handler(&rspW, &req)
	}
}

//HTTPGet 单纯的网络读取
func HTTPGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
