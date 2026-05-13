package logctx

import (
	"context"
	"log/slog"
)

type slogHandler struct {
	inner slog.Handler
}

func NewSlogHandler(inner slog.Handler) slog.Handler {
	return &slogHandler{inner: inner}
}

func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, f := range *Get(ctx) {
		record.AddAttrs(slog.Attr(f))
	}
	return h.inner.Handle(ctx, record)
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &slogHandler{inner: h.inner.WithAttrs(attrs)}
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	return &slogHandler{inner: h.inner.WithGroup(name)}
}
