package svrkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/alecthomas/template"
)

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

//WriteHTML 渲染网页 并输出渲染结果
func (rw *ResponseWriter) WriteHTML(file string, data interface{}) {
	rw.Write(rw.RenderHTML(file, data))
}

//Redirect 重定向
func (rw *ResponseWriter) Redirect(r *Request, url string, code int) {
	http.Redirect(rw, r.HTTPRequest(), url, code)
}

//Request HTTP请求，封装一些便捷操作
type Request struct {
	//用户自定义使用的业务信息，用于请求处理链上传递信息
	UserInfo map[string]interface{}

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

//HTTPRequest 返回原始 request
func (r *Request) HTTPRequest() *http.Request {
	return r.Request
}

//HTTPHandlerFunc svrkit 的 handler 定义
type HTTPHandlerFunc func(*ResponseWriter, *Request)

//HTTPFunc 包装svrkit 的HTTPHandlerFunc成标准 http.HandlerFunc
func HTTPFunc(handler HTTPHandlerFunc) http.HandlerFunc {
	return handler.ServeHTTP
}

func (h HTTPHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rspW := ResponseWriter{rw}
	req := Request{Request: r}

	if info := requestUserInfo.Get(fmt.Sprintf("%p", r)); info != nil {
		req.UserInfo = info.(map[string]interface{})
	}
	h(&rspW, &req)
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

//NoIndexFilesystem 用于 http.FileServer 用于屏蔽默认目录列表（被认为不安全） 代码来自 brad
//https://groups.google.com/forum/#!topic/golang-nuts/bStLPdIVM6w
type NoIndexFilesystem struct {
	FS http.FileSystem
}

//Open 覆盖FileSystem的 Open 方法 插入不读目录的包装
func (fs NoIndexFilesystem) Open(name string) (http.File, error) {
	f, err := fs.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
