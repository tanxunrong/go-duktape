package duktape

// #cgo CFLAGS: -std=c99 -I./
// #cgo LDFLAGS: libduktape.a -lm
// #include "go-duktape.h"
import "C"

import (
	"fmt"
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
