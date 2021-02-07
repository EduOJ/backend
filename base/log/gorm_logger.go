package log

import (
	"context"
	"errors"
	"gorm.io/gorm"
	logger2 "gorm.io/gorm/logger"
	"time"
)

// // Interface logger interface
//type Interface interface {
//	LogMode(LogLevel) Interface
//	Info(context.Context, string, ...interface{})
//	Warn(context.Context, string, ...interface{})
//	Error(context.Context, string, ...interface{})
//	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
//}

type GormLogger struct{}

func (GormLogger) LogMode(_ logger2.LogLevel) logger2.Interface {
	// do nothing
	// we use our own log level system.
	return GormLogger{}
}

func (GormLogger) Info(_ context.Context, msg string, param ...interface{}) {
	Infof(msg, param...)
}

func (GormLogger) Warn(_ context.Context, msg string, param ...interface{}) {
	Warningf(msg, param...)
}

func (GormLogger) Error(_ context.Context, msg string, param ...interface{}) {
	Errorf(msg, param...)
}

func (GormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		Errorf("%v, SQL: %s, rows: %v", err, sql, rows)
	}
}
