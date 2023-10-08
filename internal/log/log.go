package log

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func NewLog() {
	logger = logrus.New()
	// 配置 logrus 以输出 JSON 格式的日志
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	file, err := os.OpenFile("freebox.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}
	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		logger.SetOutput(fileAndStdoutWriter)
	} else {
		logger.Info("failed to log to file.")
	}

	// 设置日志级别（例如，Info、Error、Debug）
	logger.SetLevel(logrus.InfoLevel)

}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}
