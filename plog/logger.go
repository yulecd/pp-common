package plog

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
)

type Entry struct {
	*logrus.Entry
}

type Logger struct {
	logger *logrus.Logger

	filePrefix      string
	fileCreatedTime time.Time
	fileDir         string
	fileSplit       bool
	file            *os.File
}

func newLogger() *Logger {
	return &Logger{
		logger: logrus.New(),
	}
}

// createFile will make the Logger output to a file and split log file automatically if split is true
func (m *Logger) createFile(dir, prefix string, split bool) error {

	m.fileDir = dir
	m.filePrefix = prefix
	m.fileSplit = split

	if err := m.createFileAndSetOutPut(); err != nil {
		return err
	}

	if m.fileSplit {
		go m.checkSplitFile()
	}

	return nil
}

func (m *Logger) checkSplitFile() {
	t := time.NewTicker(time.Second)
	for range t.C {
		// check a new day
		if time.Now().Format(fileNameTimeFormatter) != m.fileCreatedTime.Format(fileNameTimeFormatter) {
			if err := m.createFileAndSetOutPut(); err != nil {
				m.logger.Errorf("createFileAndSetOutPut err=%s", err.Error())
			}
		}
	}
}

func (m *Logger) createFileAndSetOutPut() error {
	fileName := ""
	if m.fileSplit {
		fileName = fmt.Sprintf("%s.%s.log", path.Join(m.fileDir, m.filePrefix), time.Now().Format(fileNameTimeFormatter))
	} else {
		fileName = fmt.Sprintf("%s.log", path.Join(m.fileDir, m.filePrefix))
	}
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return err
	}

	m.fileCreatedTime = time.Now()
	m.logger.SetOutput(file)

	// close the old file
	if m.file != nil {
		m.file.Close()
	}

	m.file = file
	return nil
}

func (m *Logger) withFields(fields logrus.Fields) *Entry {
	return &Entry{
		Entry: m.logger.WithFields(fields),
	}
}

func (m *Logger) Debug(args ...interface{}) {
	m.logger.Debug(args...)
}

func (m *Logger) Debugf(format string, args ...interface{}) {
	m.logger.Debugf(format, args...)
}

func (m *Logger) Info(args ...interface{}) {
	m.logger.Info(args...)
}

func (m *Logger) Infof(format string, args ...interface{}) {
	m.logger.Infof(format, args...)
}

func (m *Logger) Warn(args ...interface{}) {
	m.logger.Warn(args...)
}

func (m *Logger) Warnf(format string, args ...interface{}) {
	m.logger.Warnf(format, args...)
}

func (m *Logger) Error(args ...interface{}) {
	m.logger.Error(args...)
}

func (m *Logger) Errorf(format string, args ...interface{}) {
	m.logger.Errorf(format, args...)
}

func (m *Logger) Fatal(args ...interface{}) {
	m.logger.Fatal(args...)
}

func (m *Logger) Fatalf(format string, args ...interface{}) {
	m.logger.Fatalf(format, args...)
}

func (m *Logger) SetLevel(level logrus.Level) {
	m.logger.SetLevel(level)
}
