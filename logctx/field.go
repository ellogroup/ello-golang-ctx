package logctx

import (
	"log/slog"
	"time"
)

type Field slog.Attr

func String(key, val string) Field {
	return Field(slog.String(key, val))
}

func Int(key string, val int) Field {
	return Field(slog.Int(key, val))
}

func Float64(key string, val float64) Field {
	return Field(slog.Float64(key, val))
}

func Bool(key string, val bool) Field {
	return Field(slog.Bool(key, val))
}

func Time(key string, val time.Time) Field {
	return Field(slog.Time(key, val))
}

func Duration(key string, val time.Duration) Field {
	return Field(slog.Duration(key, val))
}

func Any(key string, val interface{}) Field {
	return Field(slog.Any(key, val))
}
