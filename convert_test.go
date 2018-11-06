package svrkit

import "testing"

func TestPrettySize(t *testing.T) {
	t.Log(PrettySize(330776))
	t.Log(PrettySize(102400))
}
