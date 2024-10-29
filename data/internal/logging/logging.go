// internal/logging/logging.go
package logging

import (
	"os"
	"time"

	"io"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	// 日志文件路径
	logPath := "/var/log/app/crawler.log"

	// 配置rotatelogs
	writer, err := rotatelogs.New(
		logPath+".%Y%m%d",
		rotatelogs.WithLinkName(logPath),          // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(12*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		log.Fatalf("Failed to initialize log file rotator: %v", err)
	}

	// 同时输出到文件和标准输出（可选）

	log.SetOutput(io.MultiWriter(os.Stdout, writer))
	log.SetReportCaller(true)

	// 设置日志格式
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.RFC3339,
	})

	// 设置日志级别
	log.SetLevel(log.InfoLevel)
}
