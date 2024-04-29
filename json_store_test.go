package svrkit

import (
	"os"
	"testing"
)

func TestJSONStore_Load(t *testing.T) {
	type dataStuct struct{ A string }
	// a := dataStuct{"hi"}
	b := NewJSONStore[dataStuct]()
	tmp, _ := os.CreateTemp("", "test_*")
	tmp.Write([]byte(`{"A":"hi"}`))
	b.Load(tmp.Name())
	if b.Data.A != "hi" {
		t.Error("decode error")
	}
}

func TestJSONSave(t *testing.T) {
	type dataStuct struct{ A string }
	b := NewJSONStore[dataStuct]()
	b.Data.A = "hello"
	tmp, _ := os.CreateTemp("", "test_*")
	b.Filename = tmp.Name()
	b.Save()

	buf := make([]byte, 100)
	n, _ := tmp.Read(buf)

	if string(buf[:n]) != `{
    "A": "hello"
}
` {
		t.Error("unexpected content")
	}
}
