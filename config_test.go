package svrkit

import (
	"os"
	"testing"
)

func TestGetArray(t *testing.T) {
	var c Config
	c.LoadFromString("")

	t.Logf("%#+v", c.GetArray("test"))
}

func TestGetFirst(t *testing.T) {
	var c Config
	c.LoadFromString("")

	t.Logf("%#+v", c.GetFirst("test"))
}

func TestConfigLoad(t *testing.T) {
	f, _ := os.CreateTemp("", "test_*")
	f.Write([]byte(`a=b`))
	var c Config
	c.Load(f.Name())
	if c.Get("a") != "b" {
		t.Error("unexpected config read")
	}
}

func TestConfigLoadOverride(t *testing.T) {
	f, _ := os.CreateTemp("", "test_*")
	f.Write([]byte(`a=b`))

	os.WriteFile(f.Name()+".dev", []byte(`a=c`), 0755)

	var c Config
	c.Load(f.Name())
	if c.Get("a") != "c" {
		t.Error("unexpected config read")
	}
}
