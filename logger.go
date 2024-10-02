package myslog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

func init() {
	defaultLogger = slog.New(NewPrettyHandler(WithWriter(os.Stdout), WithLever(slog.LevelInfo), WithCallerDepth(2)))
}

func InitDefaultLogger(writer io.Writer, logLevel slog.Level, options ...HandlerOption) {
	options = append(options, WithWriter(writer), WithLever(logLevel), WithCallerDepth(2))
	defaultLogger = slog.New(NewPrettyHandler(options...))
}

func ContextWithArgs(ctx context.Context, kvs ...any) context.Context {
	var args []any
	if ctxKv := ctx.Value(contextArgsKey); ctxKv != nil {
		args = ctxKv.([]any)
	}
	args = append(args, kvs...)
	return context.WithValue(ctx, contextArgsKey, args)
}

func Debug(ctx context.Context, format string, args ...any) {
	printLog(ctx, slog.LevelDebug, format, args...)
}

func Info(ctx context.Context, format string, args ...any) {
	printLog(ctx, slog.LevelInfo, format, args...)
}

func Warn(ctx context.Context, format string, args ...any) {
	printLog(ctx, slog.LevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...any) {
	printLog(ctx, slog.LevelError, format, args...)
}

func printLog(ctx context.Context, level slog.Level, format string, args ...any) {
	defaultLogger.Log(ctx, level, fmt.Sprintf(format, args...))
}
