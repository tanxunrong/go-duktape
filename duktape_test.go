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
