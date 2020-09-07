package dlog

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func newLfsHook(root string, td time.Duration, maxRemainCnt uint64) (logrus.Hook, error) {
	writer, err := rotatelogs.New(
		root+".%Y%m%d%H",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		// rotatelogs.WithLinkName(globalLogName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(td),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数。
		//rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(uint(maxRemainCnt)),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
		return nil, err
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{DisableColors: true})

	return lfsHook, nil
}

func InitIO(name string, level logrus.Level, json bool, td time.Duration, count uint64) (err error) {
	if len(name) == 0 {
		return
	}

	log.Out = &NullOutput{}

	if json {
		log.Formatter = &logrus.JSONFormatter{
			TimestampFormat:  "",
			DisableTimestamp: true,
		}
	} else {
		log.Formatter = &logrus.TextFormatter{
			ForceColors:            false,
			DisableColors:          true,
			DisableTimestamp:       true,
			FullTimestamp:          false,
			TimestampFormat:        "",
			DisableSorting:         true,
			DisableLevelTruncation: false,
		}
	}

	one, err := newLfsHook(name, td, count)
	if err != nil {
		return
	}

	log.AddHook(one)
	log.SetLevel(level)

	return nil
}

func LogOpt() *logrus.Logger {
	return log
}

func AppendLog(args ...interface{}) {

	log.Info(args...)

	return
}
