package logctx

import (
	"go.uber.org/zap"
	"time"
)

type Field zap.Field

func String(key, val string) Field {
	return Field(zap.String(key, val))
}

func Int(key string, val int) Field {
	return Field(zap.Int(key, val))
}

func Float32(key string, val float32) Field {
	return Field(zap.Float32(key, val))
}

func Float64(key string, val float64) Field {
	return Field(zap.Float64(key, val))
}

func Bool(key string, val bool) Field {
	return Field(zap.Bool(key, val))
}

func Time(key string, val time.Time) Field {
	return Field(zap.Time(key, val))
}

func Duration(key string, val time.Duration) Field {
	return Field(zap.Duration(key, val))
}

func Any(key string, val interface{}) Field {
	return Field(zap.Any(key, val))
}

func Error(err error) Field {
	return Field(zap.Error(err))
}
