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

	i := 2123
	c.PushInt(i)
	if f,err := c.GetNumber(-1);err != nil {
		t.Error(err)
	} else if int(f) != i {
		t.Fatal("get != push")
	}

	pi := 3.14159265369
	c.PushDouble(pi)
	if pr,err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if pr != pi {
		t.Fatal("get float != push")
	}

	c.Close()
}

func TestStr(t *testing.T) {
	c := NewCtx()
	s := "go\U00010000go\U00013000go\u27f0"
	c.PushStr(s)
	if sr,err := c.GetStr(-1) ; err != nil {
		t.Error(err)
	} else if sr != s {
		t.Fatal("get str != push")
	}
	c.Close()
}

func TestBool(t *testing.T) {
	c := NewCtx()
	c.PushBool(true)
	if ok,err := c.GetBool(-1) ; err != nil {
		t.Error(err)
	} else if !ok {
		t.Fatal("get false when push true")
	}
	c.PushBool(false)
	if ok,err := c.GetBool(-1) ; err != nil {
		t.Error(err)
	} else if ok {
		t.Fatal("get true when push false")
	}
	c.Close()
}
