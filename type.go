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
	DUK_COMPILE_EVAL     C.duk_uint_t = 1  /* compile eval code (instead of program) */
	DUK_COMPILE_FUNCTION C.duk_uint_t = 2  /* compile function code (instead of program) */
	DUK_COMPILE_STRICT   C.duk_uint_t = 4  /* use strict (outer) context for program, eval, or function */
	DUK_COMPILE_SAFE     C.duk_uint_t = 8  /* (C.duk_uint_ternal) catch compilation errors */
	DUK_COMPILE_NORESULT C.duk_uint_t = 16 /* (C.duk_uint_ternal) omit eval result */
	DUK_COMPILE_NOSOURCE C.duk_uint_t = 32 /* (C.duk_uint_ternal) no source string on stack */
	DUK_COMPILE_STRLEN   C.duk_uint_t = 64 /* (C.duk_uint_ternal) take strlen() of src_buffer (avoids double evaluation in macro) */
)

/* Enumeration flags for duk_enum() */
const (
	DUK_ENUM_INCLUDE_NONENUMERABLE C.duk_uint_t = (1 << 0) /* enumerate non-numerable properties in addition to enumerable */
	DUK_ENUM_INCLUDE_INTERNAL      C.duk_uint_t = (1 << 1) /* enumerate internal properties  C.duk_uint_t = (regardless of enumerability) */
	DUK_ENUM_OWN_PROPERTIES_ONLY   C.duk_uint_t = (1 << 2) /* don't walk prototype chain, only check own properties */
	DUK_ENUM_ARRAY_INDICES_ONLY    C.duk_uint_t = (1 << 3) /* only enumerate array indices */
	DUK_ENUM_SORT_ARRAY_INDICES    C.duk_uint_t = (1 << 4) /* sort array indices, use with DUK_ENUM_ARRAY_INDICES_ONLY */
	DUK_ENUM_NO_PROXY_BEHAVIOR     C.duk_uint_t = (1 << 5) /* enumerate a proxy object itself without invoking proxy behavior */

)

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

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1
const DUK_INVALID_INDEX int = MinInt
