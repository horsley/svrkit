package svrkit

import "testing"

func TestGetArray(t *testing.T) {
	var c Config
	c.LoadFromString("")

	t.Logf("%#+v", c.GetArray("test"))
}
