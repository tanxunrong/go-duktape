package duktape

import (
	"testing"
)

func TestOpenAndClose(t *testing.T) {
	c := NewCtx()
	if c.ctx == nil {
		t.Fatal("ctx nil")
	}
	c.Close()
	if c.ctx != nil {
		t.Fatal("ctx close fail")
	}
}

func TestPushNumber(t *testing.T) {
	c := NewCtx()

	c.PushInt(2123)
	f := c.GetNumber(-1)
	if int(f) != 2123 {
		t.Fatal("get != push")
	}

	pi := 3.14159265369
	c.PushDouble(pi)
	if c.GetNumber(-1) != pi {
		t.Fatal("get float != push")
	}

	c.Close()
}

func TestStr(t *testing.T) {
	c := NewCtx()
	s := "go\U00010000go\U00013000go\u27f0"
	c.PushStr(s)
	if c.GetStr(-1) != s {
		t.Fatal("get str != push")
	}
}
