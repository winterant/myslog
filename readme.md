# myslog

为slog增加了一种Handler，能够打印出易于浏览的日志格式。

## 🧭 使用示例

### 使用默认logger

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

func InitLogger() {
	writers := io.MultiWriter(&lumberjack.Logger{
		Filename:   "./log/main.log", // 日志文件的位置
		MaxSize:    128,              // 文件最大大小（单位MB）
		MaxBackups: 0,                // 保留的最大旧文件数量
		MaxAge:     90,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
		LocalTime:  true,             // 使用本地时间创建时间戳
	}, os.Stdout)

	myslog.InitDefaultLogger(writers, slog.LevelDebug)
}

func main() {
	ctx := context.Background()

	InitLogger()

	ctx = myslog.ContextWithArgs(ctx, "taskId", "tsk-thisisataskid") // 利用context确保每一条都输出某些信息

	myslog.Debug(ctx, "process is starting...")

	name := "Winterant"
	myslog.Info(ctx, "My name is %s.", name)
}
```

日志：

```
2024-10-02 11:42:17.227797 DEBUG /Users/jinglong/Projects/github/myslog/main.go:33 [taskId=tsk-thisisataskid] process is starting...
2024-10-02 11:42:17.228035 INFO  /Users/jinglong/Projects/github/myslog/main.go:36 [taskId=tsk-thisisataskid] My name is Winterant.
```

### 使用原生slog.Logger

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/winterant/myslog"
)

func GetLogger() *slog.Logger {
	writers := io.MultiWriter(&lumberjack.Logger{
		Filename:   "./log/main.log", // 日志文件的位置
		MaxSize:    128,              // 文件最大大小（单位MB）
		MaxBackups: 0,                // 保留的最大旧文件数量
		MaxAge:     90,               // 保留旧文件的最大天数
		Compress:   false,            // 是否压缩/归档旧文件
		LocalTime:  true,             // 使用本地时间创建时间戳
	}, os.Stdout)

	handler := myslog.NewPrettyHandler(myslog.WithWriter(writers), myslog.WithLever(slog.LevelDebug))
	return slog.New(handler).With("key", "display_in_each_log")
}

func main() {
	ctx := context.Background()

	slogger := GetLogger()

	ctx = myslog.ContextWithArgs(ctx, "taskId", "tsk-thisisatask")

	slogger.Log(ctx, slog.LevelDebug, "process is starting...")

	name := "Winterant"
	slogger.Log(ctx, slog.LevelInfo, fmt.Sprintf("My name is %s.", name), "money", "9999999")
}
```

日志：

```
2024-10-01 21:05:59.713409 DEBUG /Users/jinglong/Projects/github/myslog/main.go:35 [key=display_in_each_log] [taskId=tsk-thisisatask] process is starting...
2024-10-01 21:05:59.714219 INFO  /Users/jinglong/Projects/github/myslog/main.go:38 [key=display_in_each_log] [taskId=tsk-thisisatask] [money=9999999] My name is Winterant.
```