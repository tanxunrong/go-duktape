package duktape

/*
#cgo CFLAGS: -std=c99 -I./
#include "duktape.h"

extern void go_duktape_fatal(duk_context *ctx, duk_errcode_t code,const char* msg);

void go_duktape_fatal_cgo(duk_context *ctx, duk_errcode_t code,const char* msg) {
	go_duktape_fatal(ctx,code,msg);
}
*/
import "C"
