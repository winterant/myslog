package myslog

import (
	"log/slog"
)

func NewSlog(options ...Option) *slog.Logger {
	return slog.New(NewPrettyHandler(options...))
}
