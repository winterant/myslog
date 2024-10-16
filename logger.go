package myslog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var defaultLogger *slog.Logger

func init() {
	writers := io.MultiWriter(&lumberjack.Logger{
		Filename:   "./log/main.log", // 日志文件的位置
		MaxSize:    128,              // 文件最大大小（单位MB）
		MaxBackups: 0,                // 保留的最大旧文件数量
		MaxAge:     90,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
		LocalTime:  true,             // 使用本地时间创建时间戳
	}, os.Stdout)
	InitDefaultLogger(writers, slog.LevelDebug)
}

// InitDefaultLogger reinitializes the default logger instead of acquiescent.
func InitDefaultLogger(writer io.Writer, logLevel slog.Level, options ...HandlerOption) {
	options = append(options, WithWriter(writer), WithLever(logLevel), WithCallerDepth(2))
	defaultLogger = slog.New(NewPrettyHandler(options...))
}

// ContextWithArgs returns a context with key-values which myslog will print.
func ContextWithArgs(ctx context.Context, kvs ...any) context.Context {
	var args []any
	if ctxKv := ctx.Value(contextArgsKey); ctxKv != nil {
		args = ctxKv.([]any)
	}
	args = append(args, kvs...)
	return context.WithValue(ctx, contextArgsKey, args)
}

func Debug(ctx context.Context, format string, args ...any) {
	log(ctx, slog.LevelDebug, format, args...)
}

func Info(ctx context.Context, format string, args ...any) {
	log(ctx, slog.LevelInfo, format, args...)
}

func Warn(ctx context.Context, format string, args ...any) {
	log(ctx, slog.LevelWarn, format, args...)
}

func Error(ctx context.Context, format string, args ...any) {
	log(ctx, slog.LevelError, format, args...)
}

func log(ctx context.Context, level slog.Level, format string, args ...any) {
	defaultLogger.Log(ctx, level, safeSprintf(format, args...))
}

func safeSprintf(format string, args ...interface{}) string {
	var buffer bytes.Buffer
	count := 0
	length := len(format)
	for i := 0; i < length; i++ {
		if format[i] == '%' {
			// Check if this '%' is independent
			isIndependent := true
			if (i > 0 && format[i-1] == '%') || (i < length-1 && format[i+1] == '%') {
				isIndependent = false
			}
			// If it's an independent '%', decide whether to keep or replace
			if isIndependent {
				if count < len(args) {
					buffer.WriteByte('%')
				} else {
					buffer.WriteString("%%")
				}
				count++
			} else {
				buffer.WriteByte('%')
			}
		} else {
			buffer.WriteByte(format[i])
		}
	}
	return fmt.Sprintf(buffer.String(), args...)
}
