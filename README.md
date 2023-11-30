# Ello Go Context packages

Common packages for handling context.

## logctx

The logctx package supports logger-specific context being added, read and modified within context.Context. This makes it 
easier to pass values such as a request id throughout the application and include it within all logs.

## logctx.Get(context.Context) *LogCtx

Gets the current `LogCtx` from `context.Context`. `LogCtx` is a slice of `logctx.Field`. If `LogCtx` is not found, a new 
`LogCtx` struct is created.

## logctx.Add(context.Context, ...Field) context.Context

Adds fields to the `LogCtx` within `context.Context`, and returns a new `context.Context`. If `LogCtx` is not found, a 
new `LogCtx` struct is created.

## logctx.Field

`logctx.Field` represents a key-value pair. There are helpers to create fields of different types, such as 
`logctx.String(key, val string) Field` and `logctx.Int(key string, val int) Field`.

## logctx.Zap(context.Context, ...zap.Field) []zap.Field

Extract `LogCtx` from `context.Context` and convert into a slice of `zap.Field` to attach to a Zap log entry.

### Example 

```go
// Example Zap logger
log := zap.NewExample()

// Example context
ctx := context.TODO()

// Add key-value pair to context
ctx = logctx.Add(ctx, logctx.Int("example_one", 123))

// Log with log context
log.Info("message", logctx.Zap(ctx)...)

// Log with log context and additional context
log.Info("message", logctx.Zap(ctx, logctx.Bool("example_two", true))...)
```