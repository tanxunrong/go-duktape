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

func TestPushAndPushArr(t *testing.T) {
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

	var u1 uint16 = 324
	var i1 int64 = 2392931293
	c.Push(u1)
	if pr, err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if pr != float64(u1) {
		t.Fatal("get float != push")
	}
	c.Push(i1)
	if pr, err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if pr != float64(i1) {
		t.Fatal("get float != push")
	}
	c.PopN(4)

	arr := []interface{}{"abc", 123, u1, i1}
	c.PushArr(arr)

	marr := make(map[string]interface{}, 5)
	marr["go"] = "abc"
	marr["do"] = 123
	marr["first"] = u1
	c.PushArr(marr)

	submarr := make(map[string]interface{}, 5)
	submarr["marr"] = marr
	submarr["arr"] = arr
	submarr["simple"] = u1
	c.PushArr(submarr)

	c.PushNull()
	c.PushUndefined()
	c.PopN(2)

	n := c.GetTop()
	if n != 3 {
		t.Fatalf("unexpected top num %d, dump %s", n, c.Dump())
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

	if c.Dump() != "ctx: top=1, stack=[\"TOUP\"]" {
		t.Fatal("unexpected dump")
	}
	c.PopN(1)

	c.Eval("\"abcd\".length")
	if ret, err := c.GetNumber(-1); err != nil {
		t.Error(err)
	} else if ret != float64(4) {
		t.Fatal("unexpected result")
	}
	c.PopN(1)

	c.Close()

}

func TestGc(t *testing.T) {
	c := NewCtx()
	defer c.Close()
	c.Gc()
}
