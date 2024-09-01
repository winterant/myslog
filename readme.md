# myslog

为slog增加了一种Handler，能够打印出易于浏览的日志格式。

## 🧭 使用示例

### 获取原生slog对象

```go
package main

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/winterant/myslog"
)

func main() {
	writers := io.MultiWriter(&lumberjack.Logger{
		Filename:   "./log/main.log", // 日志文件的位置
		MaxSize:    128,              // 文件最大大小（单位MB）
		MaxBackups: 0,                // 保留的最大旧文件数量
		MaxAge:     90,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
		LocalTime:  true,             // 使用本地时间创建时间戳
	}, os.Stdout)

	logger := myslog.NewSlog(myslog.WithWriter(writers), myslog.WithLever(slog.LevelDebug)).With("taskId", "tsk-abc12345")

	logger.Log(context.Background(), slog.LevelInfo, "This is a message.")
}
```
日志：
```
2024-09-01 12:28:52.630 INFO  /Users/admin/project/main.go:25 [taskId=tsk-abc12345] This is a message.
```