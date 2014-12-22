package duktape

/* 
#cgo CFLAGS: -std=c99 -I./
#cgo LDFLAGS: libduktape.a -lm
#include "duktape.h"

void go_duktape_fatal_cgo(duk_context *, duk_errcode_t, const char*);

*/
import "C"

import (
	"errors"
	"sync"
	"fmt"
	"io/ioutil"
	"unsafe"
)

// the duk_context holder.
type Context struct {
	ctx   unsafe.Pointer
	mutex sync.Mutex
	hell  chan DukError // fatal error chan
	dead bool
}

type CtxCenter map[unsafe.Pointer]Context

var allContext = make(CtxCenter, 32)

//export go_duktape_fatal
func go_duktape_fatal(ctx *C.duk_context, code C.duk_errcode_t, msg *C.char) {
	m := C.GoString(msg)
	d := DukError{code: code, msg: m}
	if c, ok := allContext[unsafe.Pointer(ctx)]; ok {
		c.hell <- d
		c.dead = true
	}
}

var fatal C.duk_fatal_function

func init() {
	fatal = (C.duk_fatal_function)(unsafe.Pointer(C.go_duktape_fatal_cgo))
}

// create a new duktape context.
func NewCtx() Context {
	if fatal == nil {
		panic("fatal func pointer = nil")
	}
	ctx := C.duk_create_heap(nil, nil, nil, nil, fatal)
	if ctx == nil {
		panic("new ctx = nil")
	}
	c := Context{ctx: ctx, hell: make(chan DukError, 5),dead:false}
	allContext[ctx] = c
	load,err := ioutil.ReadFile("./load.js")
	if err != nil {
		panic(err)
	}
	c.load(string(load))
	return c
}

func (c *Context) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.ctx == nil {
		panic("new ctx = nil")
	}
	delete(allContext, c.ctx)
	C.duk_destroy_heap(c.ctx)
	c.ctx = nil
	c.dead = true
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

// Push int value
func (c *Context) PushInt(i int) {
	c.check()
	C.duk_push_number(c.ctx, C.duk_double_t(float64(i)))
}

// Push float64 value
func (c *Context) PushDouble(f float64) {
	c.check()
	C.duk_push_number(c.ctx, C.duk_double_t(f))
}

// Push string value
func (c *Context) PushStr(s string) {
	c.check()
	str := C.CString(s)
	l := C.duk_size_t(len(s))
	C.duk_push_lstring(c.ctx, str, l)
}

// Push null value
func (c *Context) PushNull() {
	c.check()
	C.duk_push_null(c.ctx)
}

// Push undefined value
func (c *Context) PushUndefined() {
	c.check()
	C.duk_push_undefined(c.ctx)
}

// Push bool value
func (c *Context) PushBool(b bool) {
	c.check()
	if b {
		C.duk_push_true(c.ctx)
	} else {
		C.duk_push_false(c.ctx)
	}
}

// Get float64 value from stack index i.
// i can be -1,-2,... or n,n-1,...,1 from top to bottom.
// n == c.GetTop()
func (c *Context) GetNumber(i int) (float64, error) {
	c.check()
	b := C.duk_is_number(c.ctx, C.duk_idx_t(i))
	if b == 0 {
		return 0, TypeError
	}
	num := C.duk_get_number(c.ctx, C.duk_idx_t(i))
	return float64(num), nil
}

// Get bool value
func (c *Context) GetBool(i int) (bool, error) {
	c.check()
	b := C.duk_is_boolean(c.ctx, C.duk_idx_t(i))
	if b == 0 {
		return false, TypeError
	}
	ret := C.duk_get_boolean(c.ctx, C.duk_idx_t(i))
	if ret > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

// Get string
func (c *Context) GetStr(i int) (string, error) {
	c.check()
	b := C.duk_is_string(c.ctx, C.duk_idx_t(i))
	if b == 0 {
		return "", TypeError
	}
	var l C.int
	s := C.duk_get_lstring(c.ctx, C.duk_idx_t(i), (*C.duk_size_t)(unsafe.Pointer(&l)))
	//TODO unknown error cause panic without this line
	_ = fmt.Sprintf("%v",s)
	return C.GoStringN(s, l), nil
}

func (c *Context) GetArr(i int) (map[string]interface{},error) {
	c.check()
	b := C.duk_is_array(c.ctx, C.duk_idx_t(i))
	if b == 0 {
		return nil, TypeError
	}
	C.duk_enum(c.ctx,C.duk_idx_t(i),0)
	for {
		b = C.duk_next(c.ctx,
}

// return current number of values on stack
func (c *Context) GetTop() int {
	c.check()
	return int(C.duk_get_top(c.ctx))
}

// load string with <load> filename
func (c *Context) load(s string) {
	c.check()
	str := C.CString(s)
	l := len(s)
	c.PushStr("<load>")
	C.duk_eval_raw(c.ctx, str, (C.duk_size_t)(l),(DUK_COMPILE_EVAL | DUK_COMPILE_NOSOURCE | DUK_COMPILE_NORESULT | DUK_COMPILE_SAFE) )
	C.free(unsafe.Pointer(str))
}

// eval string with <eval> filename
func (c *Context) Eval(s string) {
	c.check()
	str := C.CString(s)
	l := len(s)
	c.PushStr("<eval>")
	C.duk_eval_raw(c.ctx, str, (C.duk_size_t)(l),(DUK_COMPILE_EVAL | DUK_COMPILE_NOSOURCE | DUK_COMPILE_SAFE) )
	C.free(unsafe.Pointer(str))
}

// Push values onto the stack.can be number,string or bool.
func (c *Context) Push(i interface{}) {
	switch i.(type) {
	case uint8:
		f := float64(i.(uint8))
		c.PushDouble(f)
		break
	case uint16:
		f := float64(i.(uint16))
		c.PushDouble(f)
		break
	case uint32:
		f := float64(i.(uint32))
		c.PushDouble(f)
		break
	case uint64:
		f := float64(i.(uint64))
		c.PushDouble(f)
		break
	case uint:
		f := float64(i.(uint))
		c.PushDouble(f)
		break
	case int8:
		f := float64(i.(int8))
		c.PushDouble(f)
		break
	case int16:
		f := float64(i.(int16))
		c.PushDouble(f)
		break
	case int32:
		f := float64(i.(int32))
		c.PushDouble(f)
		break
	case int64:
		f := float64(i.(int64))
		c.PushDouble(f)
		break
	case int:
		f := float64(i.(int))
		c.PushDouble(f)
		break
	case float32:
		f := float64(i.(float32))
		c.PushDouble(f)
		break
	case float64:
		f := i.(float64)
		c.PushDouble(f)
		break
	case string:
		s := i.(string)
		c.PushStr(s)
		break
	case bool:
		b := i.(bool)
		c.PushBool(b)
		break
	case []interface{}:
		c.check()
		idx := C.duk_push_array(c.ctx)
		arr := i.([]interface{})
		for key,val := range(arr) {
			c.Push(val)
			ret := C.duk_put_prop_index(c.ctx,idx,C.duk_uarridx_t(key));
			if int(ret) != 1 {
				panic("push failed")
			}
		}
		break
	case map[string]interface{}:
		c.check()
		idx := C.duk_push_array(c.ctx)
		arr := i.(map[string]interface{})
		for key,val := range(arr) {
			c.Push(val)
			k := C.CString(key)
			ret := C.duk_put_prop_string(c.ctx,idx,k);
			if int(ret) != 1 {
				panic("push failed")
			}
			C.free(unsafe.Pointer(k))
		}
		break
	default:
		panic("push failed,invaid type")
		break
	}
}

// pop values from stack, i >= 0.
func (c *Context) PopN(i int) {
	if i == 0 {
		return
	}
	if i < 0 {
		panic("invalid i < 0")
	}
	C.duk_pop_n(c.ctx,C.duk_idx_t(i))
}

// fatal call shall not return.Don't use it.
func (c *Context) fatal(code C.duk_errcode_t, msg string) {
	c.check()
	str := C.CString(msg)
	C.duk_fatal(c.ctx, code, str)
	C.free(unsafe.Pointer(str))
}

// check if context is dead
func (c *Context) check() {
	if c.ctx == nil || c.dead {
		if len(c.hell) > 0 {
			e := <-c.hell
			panic(e)
		}
		panic("context is dead")
	}
}

// dump the stack content
func (c *Context) Dump() string {
	C.duk_push_context_dump(c.ctx)
	var l C.duk_size_t
	s := C.duk_safe_to_lstring(c.ctx,-1,&l)
	str := C.GoStringN(s,C.int(l))
	return str
}

// run the gc
func (c *Context) Gc() {
	C.duk_gc(c.ctx,C.duk_uint_t(0))
}
