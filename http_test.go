package svrkit

import (
	"net/http"
	"testing"
)

func TestHTTPWriteString(t *testing.T) {
	http.HandleFunc("/", HTTPFunc(func(rw *ResponseWriter, r *Request) {
		rw.WriteString("hello world")
	}))
	svr := &HTTPSvr{}
	svr.Addr = ":18424"
	svr.ListenAndServe()
}
