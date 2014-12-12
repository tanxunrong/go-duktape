package duktape

// #cgo CFLAGS: -std=c99 -I./
// #include "duktape.h"
import "C"


/* Duktape specific error codes */
const (
	DUK_ERR_UNIMPLEMENTED_ERROR C.duk_errcode_t = 50 /* UnimplementedError */
	DUK_ERR_UNSUPPORTED_ERROR   C.duk_errcode_t = 51 /* UnsupportedError */
	DUK_ERR_INTERNAL_ERROR      C.duk_errcode_t = 52 /* InternalError */
	DUK_ERR_ALLOC_ERROR         C.duk_errcode_t = 53 /* AllocError */
	DUK_ERR_ASSERTION_ERROR     C.duk_errcode_t = 54 /* AssertionError */
	DUK_ERR_API_ERROR           C.duk_errcode_t = 55 /* APIError */
	DUK_ERR_UNCAUGHT_ERROR      C.duk_errcode_t = 56 /* UncaughtError */
)

/* Ecmascript E5 specification error codes */
const (
	DUK_ERR_ERROR           C.duk_errcode_t = 100 /* Error */
	DUK_ERR_EVAL_ERROR      C.duk_errcode_t = 101 /* EvalError */
	DUK_ERR_RANGE_ERROR     C.duk_errcode_t = 102 /* RangeError */
	DUK_ERR_REFERENCE_ERROR C.duk_errcode_t = 103 /* ReferenceError */
	DUK_ERR_SYNTAX_ERROR    C.duk_errcode_t = 104 /* SyntaxError */
	DUK_ERR_TYPE_ERROR      C.duk_errcode_t = 105 /* TypeError */
	DUK_ERR_URI_ERROR       C.duk_errcode_t = 106 /* URIError */
)

type DukError struct {
	code C.duk_errcode_t
	msg  string
}

/* Compilation flags for duk_compile() and duk_eval() */
const (
	 DUK_COMPILE_EVAL      C.duk_uint_t = 1         /* compile eval code (instead of program) */
	 DUK_COMPILE_FUNCTION  C.duk_uint_t = 2         /* compile function code (instead of program) */
	 DUK_COMPILE_STRICT    C.duk_uint_t = 4         /* use strict (outer) context for program, eval, or function */
	 DUK_COMPILE_SAFE      C.duk_uint_t = 8         /* (C.duk_uint_ternal) catch compilation errors */
	 DUK_COMPILE_NORESULT  C.duk_uint_t = 16         /* (C.duk_uint_ternal) omit eval result */
	 DUK_COMPILE_NOSOURCE  C.duk_uint_t = 32        /* (C.duk_uint_ternal) no source string on stack */
	 DUK_COMPILE_STRLEN    C.duk_uint_t = 64        /* (C.duk_uint_ternal) take strlen() of src_buffer (avoids double evaluation in macro) */
)
