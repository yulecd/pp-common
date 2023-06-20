package db

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/utils"
	"time"

	"github.com/yulecd/pp-common/plog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const (
	// Silent silent log level
	Silent string = "SILENT"
	// Error error log level
	Error = "ERROR"
	// Warn warn log level
	Warn = "WARN"
	// Info info log level
	Info = "INFO"
)

type MysqlConf struct {
	DataBase                  string        `yaml:"database"`
	Addr                      string        `yaml:"addr"` // host:port
	User                      string        `yaml:"user"`
	Password                  string        `yaml:"password"`
	MaxIdleConns              int           `yaml:"max_idle_conns"`
	MaxOpenConns              int           `yaml:"max_open_conns"`
	ConnMaxLifeTime           time.Duration `yaml:"conn_max_life_time"`
	ConnTimeOut               time.Duration `yaml:"conn_time_ut"`
	WriteTimeOut              time.Duration `yaml:"write_timeout"`
	ReadTimeOut               time.Duration `yaml:"read_timeout"`
	LogMode                   bool          `yaml:"log_mod"`
	LogLevel                  string        `yaml:"log_level"`
	IgnoreRecordNotFoundError bool          `yaml:"ignore_record_not_found_error"`
	SlowThreshold             time.Duration `yaml:"slow_threshold"`
}

func (conf *MysqlConf) checkConf() {
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 10
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 1000
	}
	if conf.ConnMaxLifeTime == 0 {
		conf.ConnMaxLifeTime = 3600 * time.Second
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 3 * time.Second
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 2 * time.Second
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 2 * time.Second
	}
	// sql 日志为基础的交互日志，默认都打印
	conf.LogMode = true
}

func InitMysqlClient(conf MysqlConf) (client *gorm.DB, err error) {
	conf.checkConf()

	gormOptConf := &gorm.Config{}
	if conf.LogMode {
		gormOptConf.Logger = initLogger(conf)
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=%s&readTimeout=%s&writeTimeout=%s&parseTime=True&loc=UTC",
		conf.User,
		conf.Password,
		conf.Addr,
		conf.DataBase,
		conf.ConnTimeOut,
		conf.ReadTimeOut,
		conf.WriteTimeOut)
	client, err = gorm.Open(mysql.Open(dsn), gormOptConf)
	if err != nil {
		plog.Errorf(nil, "init mysql conn error:%s", err.Error())
		return client, err
	}
	sqlDb, sqlDbErr := client.DB()
	if sqlDbErr != nil {
		plog.Errorf(nil, "init mysql handle error:%s", sqlDbErr.Error())
		return client, sqlDbErr
	}

	sqlDb.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDb.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDb.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	return client, nil
}

type GORMLogger struct {
	LogLevel                  gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	Colorful                  bool
}

func (l *GORMLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l GORMLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		return
	}

	warnStr := "mysql warn %s"
	plog.Errorf(nil, warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GORMLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		return
	}

	errStr := "mysql error %s"
	plog.Errorf(nil, errStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GORMLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		return
	}

	infoStr := "mysql info %s"
	plog.Infof(nil, infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
}

func (l GORMLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	logIns := plog.GetDefaultFieldEntry(ctx).WithField("cost", float64(elapsed.Microseconds())/1000)

	switch {
	case err != nil && l.LogLevel >= gormlogger.Error && (!l.IgnoreRecordNotFoundError || !errors.Is(err, gorm.ErrRecordNotFound)):
		sql, _ := fc()
		errStr := "%s [SQL]:%s"
		logIns.Errorf(errStr, err.Error(), sql)
	case l.SlowThreshold != 0 && elapsed > l.SlowThreshold && l.LogLevel >= gormlogger.Warn:
		sql, _ := fc()
		errStr := "slow sql warning [SQL]:%s"
		logIns.Warnf(errStr, sql)
	case l.LogLevel >= gormlogger.Info:
		sql, _ := fc()
		logIns.Infof("mysql exec [SQL]:%s", sql)
	}
}

func initLogger(conf MysqlConf) (logHandle *GORMLogger) {
	if !conf.LogMode {
		return
	}

	logHandle = &GORMLogger{}

	var logLevel gormlogger.LogLevel
	switch conf.LogLevel {
	case Silent:
		logLevel = gormlogger.Silent
	case Error:
		logLevel = gormlogger.Error
	case Warn:
		logLevel = gormlogger.Warn
	case Info:
		logLevel = gormlogger.Info
	default:
		logLevel = gormlogger.Info
	}

	logHandle.LogLevel = logLevel

	if conf.IgnoreRecordNotFoundError {
		logHandle.IgnoreRecordNotFoundError = conf.IgnoreRecordNotFoundError // 默认 false
	}

	if conf.SlowThreshold > 0 {
		logHandle.SlowThreshold = conf.SlowThreshold
	}

	return logHandle
}
