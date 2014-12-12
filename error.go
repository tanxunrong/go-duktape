package duktape

import (
	"fmt"
)

func (e *DukError) Error() string {
	var desc string
	switch e.code {
	case DUK_ERR_UNIMPLEMENTED_ERROR:
		desc = "DUK_ERR_UNIMPLEMENTED_ERROR"
		break
	case DUK_ERR_UNSUPPORTED_ERROR:
		desc = "DUK_ERR_UNSUPPORTED_ERROR"
		break
	case DUK_ERR_INTERNAL_ERROR:
		desc = "DUK_ERR_INTERNAL_ERROR"
		break
	case DUK_ERR_ALLOC_ERROR:
		desc = "DUK_ERR_ALLOC_ERROR"
		break
	case DUK_ERR_ASSERTION_ERROR:
		desc = "DUK_ERR_ASSERTION_ERROR"
		break
	case DUK_ERR_API_ERROR:
		desc = "DUK_ERR_API_ERROR"
		break
	case DUK_ERR_UNCAUGHT_ERROR:
		desc = "DUK_ERR_UNCAUGHT_ERROR"
		break
	case DUK_ERR_ERROR:
		desc = "DUK_ERR_ERROR"
		break
	case DUK_ERR_EVAL_ERROR:
		desc = "DUK_ERR_EVAL_ERROR"
		break
	case DUK_ERR_RANGE_ERROR:
		desc = "DUK_ERR_RANGE_ERROR"
		break
	case DUK_ERR_REFERENCE_ERROR:
		desc = "DUK_ERR_REFERENCE_ERROR"
		break
	case DUK_ERR_SYNTAX_ERROR:
		desc = "DUK_ERR_SYNTAX_ERROR"
		break
	case DUK_ERR_TYPE_ERROR:
		desc = "DUK_ERR_TYPE_ERROR"
		break
	case DUK_ERR_URI_ERROR:
		desc = "DUK_ERR_URI_ERROR"
		break
	default:
		desc = "DUK_ERR_ERROR"
		break
	}

	if len(e.msg) > 0 {
		return desc
	} else {
		return fmt.Sprintf("%v : %v", desc, e.msg)
	}
}
