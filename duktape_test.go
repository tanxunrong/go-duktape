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
	if f, err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if int(f) != i {
		t.Fatal("get != push")
	}

	pi := 3.14159265369
	c.PushDouble(pi)
	if pr, err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if pr != pi {
		t.Fatal("get float != push")
	}

	if count := c.GetTop(); count != 2 {
		t.Fatalf("stack values should be 2 instead of %v", c)
	}

	c.Close()
}

func TestStr(t *testing.T) {
	c := NewCtx()
	s := "go\U00010000go\U00013000go\u27f0"
	c.PushStr(s)
	if sr, err := c.GetStr(-1); err != nil {
		t.Error(err)
	} else if sr != s {
		t.Fatal("get str != push")
	}
	c.Close()
}

func TestBool(t *testing.T) {
	c := NewCtx()
	c.PushBool(true)
	if ok, err := c.GetBool(-1); err != nil {
		t.Error(err)
	} else if !ok {
		t.Fatal("get false when push true")
	}
	c.PushBool(false)
	if ok, err := c.GetBool(-1); err != nil {
		t.Error(err)
	} else if ok {
		t.Fatal("get true when push false")
	}
	c.Close()
}

/*
func TestFatalCallback(t *testing.T) {
	c := NewCtx()
	c.fatal(DUK_ERR_ERROR, "fatal error for callback test")
	if len(c.hell) < 1 {
		t.Fatal("expect fatal error")
	}
	e := <-c.hell
	t.Logf("the fatal: %v", e)
	c.Close()
}
*/

func TestEvalAndDump(t *testing.T) {
	c := NewCtx()
	c.Eval("\"toup\".toUpperCase()")
	if ret, err := c.GetStr(-1); err != nil {
		t.Error(err)
	} else if ret != "TOUP" {
		t.Fatal("unexpected result")
	}

	if c.dump() != "ctx: top=1, stack=[\"TOUP\"]" {
		t.Fatal("unexpected dump")
	}
	c.Close()

	/*
	c.Eval("\"abcd\".length")
	if ret, err := c.GetStr(-1); err != nil {
		t.Errorf()
	} else if ret != "4" {
		t.Fatal("unexpected result")
	}
	*/

}
