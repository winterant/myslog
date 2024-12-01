package myslog

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	writers := io.MultiWriter(&lumberjack.Logger{
		Filename:   "./log/main.log", // log file path
		MaxSize:    128,              // file max size in MB
		MaxBackups: 0,                // max number of backup log files
		MaxAge:     90,               // max number of days to keep old files
		Compress:   false,            // whether to compress/archive old files
		LocalTime:  true,             // Use local time or not
	}, os.Stdout)
	InitDefaultLogger(writers, slog.LevelDebug)
}
