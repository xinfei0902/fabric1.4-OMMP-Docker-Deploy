package dlog

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	log = logrus.New()
}

// Global Default Values
var (
	GlobalConsoleMark  = "console"
	GlobalDebugLevel   = logrus.DebugLevel
	GlobalDefaultLevel = logrus.InfoLevel
)

// InitLog to start log
func InitLog(name string, level string, json bool, td time.Duration, count uint64) (err error) {

	var base = logrus.InfoLevel
	switch level {
	case logrus.DebugLevel.String():
		base = logrus.DebugLevel
	case logrus.InfoLevel.String():
		base = logrus.InfoLevel
	case logrus.WarnLevel.String():
		base = logrus.WarnLevel
	case logrus.ErrorLevel.String():
		base = logrus.ErrorLevel
	case logrus.FatalLevel.String():
		base = logrus.FatalLevel
	default:
		base = GlobalDefaultLevel
	}

	logrus.SetLevel(base)

	if name == GlobalConsoleMark {
		InitConsole(base, json)
		return
	}

	if td < 1*time.Hour {
		td = 24 * time.Hour
	}

	if count <= 0 {
		count = 7
	}

	return InitIO(name, base, json, td, count)
}

// InitConsole log output into console
func InitConsole(level logrus.Level, json bool) {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//标准输出指向io.Writer固定输出位置
	log.Out = os.Stdout
	setFormat(level, json)
}

//设置日志格式 json或者文本
func setFormat(level logrus.Level, json bool) {
	log.SetLevel(level)

	if json {
		// Log as JSON instead of the default ASCII formatter.
		log.Formatter = &logrus.JSONFormatter{}
	} else {
		log.Formatter = &logrus.TextFormatter{}
	}
}

func DebugLog(name string, serial uint64) *logrus.Entry {
	return log.WithField("model", name).WithField("serial", serial)
}

func Debug(msg string) {
	log.Debug(msg)
}

func Debugf(format string, args ...interface{}) {
	log.Debug(fmt.Sprintf(format, args...))
}

func Info(msg string) {
	log.Info(msg)
}

func Warn(msg string) {
	log.Warn(msg)
}

func Error(err error) {
	log.Error(err)
}

func Fatal(msg string) {
	log.Fatal(msg)
}
