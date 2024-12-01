package myslog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

var contextArgsKey int

// PrettyHandler implements a slog handler with pretty format as easy to view.
type PrettyHandler struct {
	slog.Handler
	addSource   bool
	level       slog.Level
	callerDepth int
	w           io.Writer
	logAttrs    []slog.Attr
	mu          *sync.Mutex
}

type HandlerOption func(*PrettyHandler)

func WithWriter(writer io.Writer) HandlerOption {
	return func(handler *PrettyHandler) {
		handler.w = writer
	}
}

func WithCodeSource(addSource bool) HandlerOption {
	return func(handler *PrettyHandler) {
		handler.addSource = addSource
	}
}

func WithLever(level slog.Level) HandlerOption {
	return func(handler *PrettyHandler) {
		handler.level = level
	}
}

func WithCallerDepth(depth int) HandlerOption {
	return func(handler *PrettyHandler) {
		handler.callerDepth = depth
	}
}

func NewPrettyHandler(options ...HandlerOption) *PrettyHandler {
	handler := PrettyHandler{
		addSource:   true,
		level:       slog.LevelInfo,
		callerDepth: 0,
		w:           os.Stdout,
		mu:          &sync.Mutex{},
	}
	for _, option := range options {
		option(&handler)
	}
	handler.Handler = slog.NewJSONHandler(handler.w, &slog.HandlerOptions{AddSource: handler.addSource, Level: handler.level})
	return &handler
}

func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Handler.Enabled(ctx, level)
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.logAttrs = append(h.logAttrs, attrs...)
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	h.Handler = h.Handler.WithGroup(name)
	return h
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	buf.WriteString(fmt.Sprintf("%s %-5s", r.Time.Format("2006-01-02 15:04:05.000000"), r.Level.String()))

	if h.addSource {
		if _, file, line, ok := runtime.Caller(3 + h.callerDepth); ok {
			buf.WriteString(fmt.Sprintf(" %s:%d", file, line))
		}
	}

	for _, attr := range h.logAttrs { // 创建slog.Logger时添加的参数
		buf.WriteString(fmt.Sprintf(" [%s=%s]", attr.Key, attr.Value.String()))
	}
	ctxArgs := getContextArgs(ctx) // context中的参数
	for i := 0; i+1 < len(ctxArgs); i += 2 {
		buf.WriteString(fmt.Sprintf(" [%s=%s]", ctxArgs[i], ctxArgs[i+1]))
	}
	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool { // 打印日志时临时添加的参数
			buf.WriteString(fmt.Sprintf(" [%s=%s]", attr.Key, attr.Value.String()))
			return true
		})
	}

	buf.WriteString(fmt.Sprintf(" %s\n", r.Message))

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}

func getContextArgs(ctx context.Context) []any {
	v := ctx.Value(&contextArgsKey)
	if v != nil {
		return v.([]any)
	}
	return nil
}
