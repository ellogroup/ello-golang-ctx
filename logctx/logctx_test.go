package logctx

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func fieldsToLogCtx(s ...Field) *LogCtx {
	lc := LogCtx(s)
	return &lc
}

func TestGet(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *LogCtx
	}{
		{
			name: "Context does not contain LogCtx, returns new LogCtx",
			args: args{ctx: context.Background()},
			want: fieldsToLogCtx(),
		},
		{
			name: "Context contains LogCtx, returns LogCtx",
			args: args{ctx: context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("example", "test")))},
			want: fieldsToLogCtx(String("example", "test")),
		},
		{
			name: "Context contains multiple LogCtx, returns latest LogCtx",
			args: args{
				ctx: context.WithValue(
					context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("old", "test-old"))),
					ctxKey,
					fieldsToLogCtx(String("latest", "test-latest")),
				),
			},
			want: fieldsToLogCtx(String("latest", "test-latest")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Get(tt.args.ctx), "Get(%v)", tt.args.ctx)
		})
	}
}

func TestAdd(t *testing.T) {
	type args struct {
		ctx context.Context
		f   []Field
	}
	tests := []struct {
		name string
		args args
		want *LogCtx
	}{
		{
			name: "Context does not contain LogCtx and fields provided, returns LogCtx",
			args: args{
				ctx: context.Background(),
				f:   []Field{String("example-string", "test"), Int("example-int", 123)},
			},
			want: fieldsToLogCtx(String("example-string", "test"), Int("example-int", 123)),
		},
		{
			name: "Context contains LogCtx and fields provided, returns updated LogCtx",
			args: args{
				ctx: context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("example-existing", "test"))),
				f:   []Field{String("example-string", "test"), Int("example-int", 123)},
			},
			want: fieldsToLogCtx(String("example-existing", "test"), String("example-string", "test"), Int("example-int", 123)),
		},
		{
			name: "Context does not contain LogCtx and no fields provided, returns empty LogCtx",
			args: args{
				ctx: context.Background(),
				f:   []Field{},
			},
			want: fieldsToLogCtx(),
		},
		{
			name: "Context contains LogCtx and no fields provided, returns existing LogCtx",
			args: args{
				ctx: context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("example-existing", "test"))),
				f:   []Field{},
			},
			want: fieldsToLogCtx(String("example-existing", "test")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.args.ctx, tt.args.f...).Value(ctxKey).(*LogCtx)
			assert.Equalf(t, tt.want, got, "Add(%v, %v)", tt.args.ctx, tt.args.f)
		})
	}
}
