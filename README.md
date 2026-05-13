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

## logctx.NewSlogHandler(inner slog.Handler) slog.Handler

Wraps an existing `slog.Handler` so that any `LogCtx` fields stored in the `context.Context` are automatically injected
into every log record. Use this when constructing your `slog.Logger`.

### Example

```go
// Wrap your chosen handler so LogCtx fields are injected automatically
logger := slog.New(logctx.NewSlogHandler(slog.NewJSONHandler(os.Stdout, nil)))

// Example context
ctx := context.TODO()

// Add key-value pair to context
ctx = logctx.Add(ctx, logctx.Int("example_one", 123))

// Log with log context — "example_one" is included automatically
logger.InfoContext(ctx, "message")

// Additional fields can still be passed directly to the log call
logger.InfoContext(ctx, "message", slog.Bool("example_two", true))
```
