package duktape

// #cgo CFLAGS: -std=c99 -I./
// #cgo LDFLAGS: libduktape.a -lm
// #include "go-duktape.h"
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

// the duk_context holder.
type Context struct {
	ctx   *C.duk_context
	mutex sync.Mutex
	hell  chan DukError // fatal error chan
	dead bool
}

type CtxCenter map[*C.duk_context]Context

var allContext = make(CtxCenter, 32)

//export go_duktape_fatal
func go_duktape_fatal(ctx *C.duk_context, code C.duk_errcode_t, msg *C.char) {
	m := C.GoString(msg)
	d := DukError{code: code, msg: m}
	if c, ok := allContext[ctx]; ok {
		c.hell <- d
		c.dead = true
	}
}

var GoFatalCall = go_duktape_fatal

// create a new duktape context.
func NewCtx() Context {
	var ctx *C.duk_context
	var fatal C.duk_fatal_function
	fatal = (C.duk_fatal_function)(unsafe.Pointer(&GoFatalCall))
	ctx = (*C.duk_context)(C.duk_create_heap(nil, nil, nil, nil, fatal))
	c := Context{ctx: ctx, hell: make(chan DukError, 5),dead:false}
	allContext[ctx] = c
	return c
}

func (c *Context) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(allContext, c.ctx)
	C.duk_destroy_heap(unsafe.Pointer(c.ctx))
	c.ctx = nil
}

type DukType int

/* Value types, used by e.g. duk_get_type() */
const (
	DUK_TYPE_NONE      DukType = 0 /* no value, e.g. invalid index */
	DUK_TYPE_UNDEFINED DukType = 1 /* Ecmascript undefined */
	DUK_TYPE_NULL      DukType = 2 /* Ecmascript null */
	DUK_TYPE_BOOLEAN   DukType = 3 /* Ecmascript boolean: 0 or 1 */
	DUK_TYPE_NUMBER    DukType = 4 /* Ecmascript number: double */
	DUK_TYPE_STRING    DukType = 5 /* Ecmascript string: CESU-8 / extended UTF-8 encoded */
	DUK_TYPE_OBJECT    DukType = 6 /* Ecmascript object: includes objects, arrays, functions, threads */
	DUK_TYPE_BUFFER    DukType = 7 /* fixed or dynamic, garbage collected byte buffer */
	DUK_TYPE_POINTER   DukType = 8 /* raw void pointer */
)

var TypeError = errors.New("unexpected type")

func (c *Context) PushInt(i int) {
	c.check()
	C.duk_push_number(unsafe.Pointer(c.ctx), C.duk_double_t(float64(i)))
}

func (c *Context) PushDouble(f float64) {
	c.check()
	C.duk_push_number(unsafe.Pointer(c.ctx), C.duk_double_t(f))
}

func (c *Context) PushStr(s string) {
	c.check()
	str := C.CString(s)
	l := C.duk_size_t(len(s))
	C.duk_push_lstring(unsafe.Pointer(c.ctx), str, l)
}

func (c *Context) PushBool(b bool) {
	c.check()
	if b {
		C.duk_push_true(unsafe.Pointer(c.ctx))
	} else {
		C.duk_push_false(unsafe.Pointer(c.ctx))
	}
}

func (c *Context) GetNumber(i int) (float64, error) {
	c.check()
	b := C.duk_is_number(unsafe.Pointer(c.ctx), C.duk_idx_t(i))
	if b == 0 {
		return 0, TypeError
	}
	num := C.duk_get_number(unsafe.Pointer(c.ctx), C.duk_idx_t(i))
	return float64(num), nil
}

func (c *Context) GetBool(i int) (bool, error) {
	c.check()
	b := C.duk_is_boolean(unsafe.Pointer(c.ctx), C.duk_idx_t(i))
	if b == 0 {
		return false, TypeError
	}
	ret := C.duk_get_boolean(unsafe.Pointer(c.ctx), C.duk_idx_t(i))
	if ret > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (c *Context) GetStr(i int) (string, error) {
	c.check()
	b := C.duk_is_string(unsafe.Pointer(c.ctx), C.duk_idx_t(i))
	if b == 0 {
		return "", TypeError
	}
	var l C.int
	s := C.duk_get_lstring(unsafe.Pointer(c.ctx), C.duk_idx_t(i), (*C.duk_size_t)(unsafe.Pointer(&l)))
	return C.GoStringN(s, l), nil
}

// return current number of values on stack
func (c *Context) GetTop() int {
	c.check()
	return int(C.duk_get_top(unsafe.Pointer(c.ctx)))
}

func (c *Context) Eval(s string) {
	c.check()
	str := C.CString(s)
	l := len(s)
	c.PushStr("<eval>")
	C.duk_eval_raw(unsafe.Pointer(c.ctx), str, (C.duk_size_t)(l),(DUK_COMPILE_EVAL | DUK_COMPILE_NOSOURCE | DUK_COMPILE_SAFE) )
	C.free(unsafe.Pointer(str))
}

func (c *Context) fatal(code C.duk_errcode_t, msg string) {
	c.check()
	str := C.CString(msg)
	C.duk_fatal(unsafe.Pointer(c.ctx), code, str)
	C.free(unsafe.Pointer(str))
}

func (c *Context) check() {
	if c.dead {
		if len(c.hell) > 0 {
			e := <-c.hell
			panic(e)
		}
		panic("context is dead")
	}
}
