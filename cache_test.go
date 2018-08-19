package svrkit

import (
	"testing"
	"time"
)

func TestTTLCache(t *testing.T) {
	c := NewTTLCache()
	c.SetWithExpire("a", "b", time.Now().Add(1*time.Second))
	t.Log(c.GetString("a"))
	time.Sleep(2*time.Second)
	t.Log(c.GetString("a"))
}
