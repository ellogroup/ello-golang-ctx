package logctx

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestZap(t *testing.T) {
	type args struct {
		ctx context.Context
		f   []zap.Field
	}
	tests := []struct {
		name string
		args args
		want []zap.Field
	}{
		{
			name: "Context does not contain LogCtx and fields provided, returns provided fields",
			args: args{
				ctx: context.Background(),
				f:   []zap.Field{zap.String("example-string", "test"), zap.Int("example-int", 123)},
			},
			want: []zap.Field{zap.String("example-string", "test"), zap.Int("example-int", 123)},
		},
		{
			name: "Context contains LogCtx and fields provided, returns LofCtx combined with provided fields",
			args: args{
				ctx: context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("example-existing", "test"))),
				f:   []zap.Field{zap.String("example-string", "test"), zap.Int("example-int", 123)},
			},
			want: []zap.Field{zap.String("example-existing", "test"), zap.String("example-string", "test"), zap.Int("example-int", 123)},
		},
		{
			name: "Context does not contain LogCtx and no fields provided, returns nil",
			args: args{
				ctx: context.Background(),
				f:   []zap.Field{},
			},
			want: nil,
		},
		{
			name: "Context contains LogCtx and no fields provided, returns LogCtx only",
			args: args{
				ctx: context.WithValue(context.Background(), ctxKey, fieldsToLogCtx(String("example-existing", "test"))),
				f:   []zap.Field{},
			},
			want: []zap.Field{zap.String("example-existing", "test")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Zap(tt.args.ctx, tt.args.f...), "Zap(%v, %v)", tt.args.ctx, tt.args.f)
		})
	}
}
