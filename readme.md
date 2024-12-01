# myslog

基于slog增加了一种Handler，能够打印出易于浏览的日志格式。

# 🔨 安装

```bash
go get -u github.com/winterant/myslog
```

## 🪤 示例

### 使用默认logger

日志将输出到`os.Stdout`和文件`./log/main.log`。

```go
package main

import (
	"context"

	"github.com/winterant/myslog"
)

func main() {
	ctx := myslog.ContextWithArgs(context.Background(), "taskId", "tsk-thisisataskid", "tag", "mytag") // 利用context确保每一条都输出某些信息
	myslog.Debug(ctx, "(acquiescent myslog.Logger)process is starting...")
	myslog.Info(ctx, "My name is %s.", "Winterant")
}
```

日志：

```
2024-10-02 12:21:32.365340 DEBUG /Users/jinglong/Projects/github/myslog/main.go:12 [taskId=tsk-thisisataskid] [tag=mytag] process is starting...
2024-10-02 12:21:32.365816 INFO  /Users/jinglong/Projects/github/myslog/main.go:15 [taskId=tsk-thisisataskid] [tag=mytag] My name is Winterant.
```

### 手动初始化默认logger

只需添加一个init函数

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

func init() {
	// 自行指定日志输出目标
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
	ctx = myslog.ContextWithArgs(ctx, "taskId", "tsk-thisisataskid", "tag", "mytag") // 利用context确保每一条都输出某些信息
	myslog.Debug(ctx, "process is starting...")
	myslog.Info(ctx, "My name is %s.", "Winterant")
}
```

日志：

```
2024-10-02 11:42:17.227797 DEBUG /Users/jinglong/Projects/github/myslog/main.go:34 [taskId=tsk-thisisataskid] [tag=mytag] process is starting...
2024-10-02 11:42:17.228035 INFO  /Users/jinglong/Projects/github/myslog/main.go:37 [taskId=tsk-thisisataskid] [tag=mytag] My name is Winterant.
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
	slogger := GetLogger()
	ctx := myslog.ContextWithArgs(context.Background(), "taskId", "tsk-thisisatask")
	slogger.Log(ctx, slog.LevelDebug, "process is starting...")
	slogger.Log(ctx, slog.LevelInfo, fmt.Sprintf("My name is %s.", "Winterant"), "money", "9999999")
}

```

日志：

```
2024-10-01 21:05:59.713409 DEBUG /Users/jinglong/Projects/github/myslog/main.go:35 [key=display_in_each_log] [taskId=tsk-thisisatask] process is starting...
2024-10-01 21:05:59.714219 INFO  /Users/jinglong/Projects/github/myslog/main.go:38 [key=display_in_each_log] [taskId=tsk-thisisatask] [money=9999999] My name is Winterant.
```

## 🚛 附录

### filebeat日志收集配置

`filebeat.yaml`:

```yaml
filebeat.inputs:
  - type: log
    paths:
      - './log/*.log'
    multiline.pattern: '^\d{4}-\d{2}-\d{2}'
    multiline.negate: true
    multiline.match: after

processors:
  - drop_event:
      when:
        regexp:
          message: 'FILEBEAT_EXCLUDE'  # 排除包含FILEBEAT_EXCLUDE的日志

output.elasticsearch:
  hosts: [ "10.10.10.10:8200" ]
  username: "myusername"
  password: "mypassword"
  index: "my_project_online"

setup.ilm.enabled: false
setup.template.name: "my_project_online"
setup.template.pattern: "my_project_online-*"
setup.template.overwrite: false
```
