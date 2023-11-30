package logctx

import (
	"context"
	"go.uber.org/zap"
)

func Zap(ctx context.Context, f ...zap.Field) []zap.Field {
	var zf []zap.Field
	for _, lcf := range *Get(ctx) {
		zf = append(zf, zap.Field(lcf))
	}
	return append(zf, f...)
}
