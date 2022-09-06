package log

import "context"

type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

func WithContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func FromContext(ctx context.Context) Logger {
	if v, ok := ctx.Value(loggerKey).(Logger); ok {
		return v
	}
	return NewNilLogger()
}
