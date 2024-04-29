package svrkit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWriteString(t *testing.T) {
	http.HandleFunc("/", HTTPFunc(func(rw *ResponseWriter, r *Request) {
		rw.WriteString("hello world")
	}))

	rec := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Body.String() != "hello world" {
		t.Error("resp not match")
	}
}
