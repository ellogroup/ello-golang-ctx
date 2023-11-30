package logctx

import (
	"context"
)

type key struct{}

var ctxKey = &key{}

type LogCtx []Field

func Get(ctx context.Context) *LogCtx {
	if lc, ok := ctx.Value(ctxKey).(*LogCtx); ok {
		return lc
	}
	return new(LogCtx)
}

func Add(ctx context.Context, f ...Field) context.Context {
	lc := append(*Get(ctx), f...)
	return context.WithValue(ctx, ctxKey, &lc)
}
