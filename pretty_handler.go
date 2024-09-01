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

type Option func(*PrettyHandler)

func WithWriter(writer io.Writer) Option {
	return func(handler *PrettyHandler) {
		handler.w = writer
	}
}

func WithCodeSource(addSource bool) Option {
	return func(handler *PrettyHandler) {
		handler.addSource = addSource
	}
}

func WithLever(level slog.Level) Option {
	return func(handler *PrettyHandler) {
		handler.level = level
	}
}

func WithCallerDepth(depth int) Option {
	return func(handler *PrettyHandler) {
		handler.callerDepth = depth
	}
}

func NewPrettyHandler(options ...Option) *PrettyHandler {
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
	buf.WriteString(fmt.Sprintf("%s %-5s", r.Time.Format("2006-01-02 15:04:05.000"), r.Level.String()))

	if h.addSource {
		if _, file, line, ok := runtime.Caller(3 + h.callerDepth); ok {
			buf.WriteString(fmt.Sprintf(" %s:%d", file, line))
		}
	}

	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool {
			buf.WriteString(fmt.Sprintf(" [%s=%s]", attr.Key, attr.Value.String()))
			return true
		})
	}
	for _, attr := range h.logAttrs {
		buf.WriteString(fmt.Sprintf(" [%s=%s]", attr.Key, attr.Value.String()))
	}

	buf.WriteString(fmt.Sprintf(" %s\n", r.Message))

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.w.Write(buf.Bytes())
	return err
}
