package plog

import (
	"fmt"
	"os"

	"github.com/yulecd/pp-common/trace"
	"github.com/yulecd/pp-common/util"

	"github.com/sirupsen/logrus"
)

const (
	logTimeFormatter      = "2006-01-02 15:04:05.000"
	fileNameTimeFormatter = "2006-01-02"
)

const (
	PanicLevel logrus.Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

var stdLogger = newLogger()

// InitWithPath init the log module, this should be called at the very beginning of the whole program.
// dir is the dir path to store log file, prefix is log file prefix
// for example InitWithPath("/var", "prod") will create "/var/prod.2021-02-02.log"
func InitWithPath(dir, prefix string) {
	initLog(dir, prefix, true, logrus.InfoLevel)
}

// Init is same as InitWithPath but with a default dir "./"
func Init(prefix string) {
	initLog("./", prefix, true, logrus.InfoLevel)
}

func initLog(dir, prefix string, split bool, level logrus.Level) {
	// 初始化依赖
	trace.InitGenerator()

	// 等级
	stdLogger.logger.SetLevel(level)

	// 格式
	customFormatter := new(logrus.JSONFormatter)
	customFormatter.TimestampFormat = logTimeFormatter
	stdLogger.logger.SetFormatter(customFormatter)

	// 检查日志目录
	if !util.DirExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(fmt.Sprintf("initLog creat dir fail: %s", err))
		}
	}

	// 创建文件 如果有切分启动切分检查
	if err := stdLogger.createFile(dir, prefix, split); err != nil {
		panic(err)
	}

	// 标准输出
	stdLogger.logger.AddHook(&stdoutHook{})
}

// NewLogger return a complete new logger.
func NewLogger(dir, prefix string) (*Logger, error) {
	logger := newLogger()
	logger.logger.SetLevel(logrus.InfoLevel)

	customFormatter := new(logrus.JSONFormatter)
	customFormatter.TimestampFormat = logTimeFormatter
	logger.logger.SetFormatter(customFormatter)

	if err := logger.createFile(dir, prefix, true); err != nil {
		return nil, err
	}

	return logger, nil
}

// StdLogger 返回全局logger
func StdLogger() *Logger {
	return stdLogger
}

// SetLevel can change the output level of stdLogger.
// Init will set InfoLevel as default.
// example:
// 		import //common/log
// 		plog.Init("prod")
// 		plog.SetLevel(plog.WarnLevel)
func SetLevel(level logrus.Level) {
	stdLogger.logger.SetLevel(level)
}

// stdoutHook will print log to stdout
type stdoutHook struct{}

func (hook *stdoutHook) Fire(entry *logrus.Entry) error {
	// print to stdout
	content, err := entry.String()
	if err != nil {
		return err
	}

	fmt.Print(content)
	return nil
}

func (hook *stdoutHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
