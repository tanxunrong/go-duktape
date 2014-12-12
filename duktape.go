package duktape

// #cgo CFLAGS: -std=c99 -I./
// #cgo LDFLAGS: libduktape.a -lm
// #include "go-duktape.h"
import "C"

import (
	"fmt"
	"errors"
	"unsafe"
	"sync"
)

// duktape error code
type DukCode int

/* Duktape specific error codes */
const (
	 DUK_ERR_UNIMPLEMENTED_ERROR   DukCode =   50   /* UnimplementedError */
	 DUK_ERR_UNSUPPORTED_ERROR     DukCode =   51   /* UnsupportedError */
	 DUK_ERR_INTERNAL_ERROR        DukCode =   52   /* InternalError */
	 DUK_ERR_ALLOC_ERROR           DukCode =   53   /* AllocError */
	 DUK_ERR_ASSERTION_ERROR       DukCode =   54   /* AssertionError */
	 DUK_ERR_API_ERROR             DukCode =   55   /* APIError */
	 DUK_ERR_UNCAUGHT_ERROR        DukCode =   56   /* UncaughtError */
)

/* Ecmascript E5 specification error codes */
const (
	 DUK_ERR_ERROR                 DukCode =   100  /* Error */
	 DUK_ERR_EVAL_ERROR            DukCode =   101  /* EvalError */
	 DUK_ERR_RANGE_ERROR           DukCode =   102  /* RangeError */
	 DUK_ERR_REFERENCE_ERROR       DukCode =   103  /* ReferenceError */
	 DUK_ERR_SYNTAX_ERROR          DukCode =   104  /* SyntaxError */
	 DUK_ERR_TYPE_ERROR            DukCode =   105  /* TypeError */
	 DUK_ERR_URI_ERROR             DukCode =   106  /* URIError */
)

type Context struct {
	ctx *C.duk_context
	mutex sync.Mutex
}

type DukError struct {
	code DukCode
	msg string
}

func (e *DukError) Error() string {
	var desc string
	switch e.code {
	case DUK_ERR_UNIMPLEMENTED_ERROR:
		desc ="DUK_ERR_UNIMPLEMENTED_ERROR"
		break
	case DUK_ERR_UNSUPPORTED_ERROR:
		desc ="DUK_ERR_UNSUPPORTED_ERROR"
		break
	case DUK_ERR_INTERNAL_ERROR:
		desc ="DUK_ERR_INTERNAL_ERROR"
		break
	case DUK_ERR_ALLOC_ERROR:
		desc ="DUK_ERR_ALLOC_ERROR"
		break
	case DUK_ERR_ASSERTION_ERROR:
		desc ="DUK_ERR_ASSERTION_ERROR"
		break
	case DUK_ERR_API_ERROR:
		desc ="DUK_ERR_API_ERROR"
		break
	case DUK_ERR_UNCAUGHT_ERROR:
		desc ="DUK_ERR_UNCAUGHT_ERROR"
		break
	case DUK_ERR_ERROR:
		desc ="DUK_ERR_ERROR"
		break
	case DUK_ERR_EVAL_ERROR:
		desc ="DUK_ERR_EVAL_ERROR"
		break
	case DUK_ERR_RANGE_ERROR:
		desc ="DUK_ERR_RANGE_ERROR"
		break
	case DUK_ERR_REFERENCE_ERROR:
		desc ="DUK_ERR_REFERENCE_ERROR"
		break
	case DUK_ERR_SYNTAX_ERROR:
		desc ="DUK_ERR_SYNTAX_ERROR"
		break
	case DUK_ERR_TYPE_ERROR:
		desc ="DUK_ERR_TYPE_ERROR"
		break
	case DUK_ERR_URI_ERROR:
		desc ="DUK_ERR_URI_ERROR"
		break
	default:
		desc = "DUK_ERR_ERROR"
		break
	}

	if len(e.msg) > 0 {
		return desc
	} else {
		return fmt.Sprintf("%v : %v",desc,e.msg)
	}
}

//export go_duktape_fatal
func go_duktape_fatal (ctx *C.duk_context,code C.duk_errcode_t,msg *C.char) {
	m := C.GoString(msg)
	d := DukError{code:DukCode(code),msg:m}
	panic(d)
}

var GoFatalCall = go_duktape_fatal

// create a new duktape context.
func NewCtx() Context {
	var ctx *C.duk_context
	var fatal C.duk_fatal_function
	fatal = (C.duk_fatal_function) (unsafe.Pointer(&GoFatalCall))
	ctx = (*C.duk_context)(C.duk_create_heap(nil,nil,nil,nil,fatal))
	return Context{ctx: ctx}
}

func (c *Context) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	C.duk_destroy_heap(unsafe.Pointer(c.ctx))
	c.ctx = nil
}

type DukType int

/* Value types, used by e.g. duk_get_type() */
const (
	 DUK_TYPE_NONE        DukType =    0    /* no value, e.g. invalid index */
	 DUK_TYPE_UNDEFINED   DukType =    1    /* Ecmascript undefined */
	 DUK_TYPE_NULL        DukType =    2    /* Ecmascript null */
	 DUK_TYPE_BOOLEAN     DukType =    3    /* Ecmascript boolean: 0 or 1 */
	 DUK_TYPE_NUMBER      DukType =    4    /* Ecmascript number: double */
	 DUK_TYPE_STRING      DukType =    5    /* Ecmascript string: CESU-8 / extended UTF-8 encoded */
	 DUK_TYPE_OBJECT      DukType =    6    /* Ecmascript object: includes objects, arrays, functions, threads */
	 DUK_TYPE_BUFFER      DukType =    7    /* fixed or dynamic, garbage collected byte buffer */
	 DUK_TYPE_POINTER     DukType =    8    /* raw void pointer */
)

var TypeError = errors.New("unexpected type")

func (c *Context) PushInt(i int) {
	C.duk_push_number(unsafe.Pointer(c.ctx),C.duk_double_t(float64(i)))
}

func (c *Context) PushDouble(f float64) {
	C.duk_push_number(unsafe.Pointer(c.ctx),C.duk_double_t(f))
}

func (c *Context) PushStr(s string) {
	str := C.CString(s)
	l := C.duk_size_t(len(s))
	C.duk_push_lstring(unsafe.Pointer(c.ctx),str,l)
}

func (c *Context) PushBool(b bool) {
	if b {
		C.duk_push_true(unsafe.Pointer(c.ctx))
	} else {
		C.duk_push_false(unsafe.Pointer(c.ctx))
	}
}

func (c *Context) GetNumber(i int) (float64,error) {
	b := C.duk_is_number(unsafe.Pointer(c.ctx),C.duk_idx_t(i))
	if b == 0 {
		return 0,TypeError
	}
	num := C.duk_get_number(unsafe.Pointer(c.ctx),C.duk_idx_t(i))
	return float64(num),nil
}

func (c *Context) GetBool(i int) (bool,error) {
	b := C.duk_is_boolean(unsafe.Pointer(c.ctx),C.duk_idx_t(i))
	if b == 0 {
		return false,TypeError
	}
	ret := C.duk_get_boolean(unsafe.Pointer(c.ctx),C.duk_idx_t(i))
	if ret > 0 {
		return true,nil
	} else {
		return false,nil
	}
}

func (c *Context) GetStr(i int) (string,error) {
	b := C.duk_is_string(unsafe.Pointer(c.ctx),C.duk_idx_t(i))
	if b == 0 {
		return "",TypeError
	}
	var l C.int
	s := C.duk_get_lstring(unsafe.Pointer(c.ctx),C.duk_idx_t(i),(*C.duk_size_t)(unsafe.Pointer(&l)))
	return C.GoStringN(s,l),nil
}

// return current number of values on stack
func (c *Context) GetTop() int {
	return int(c.duk_get_top(c.ctx))
}

