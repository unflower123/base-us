package custom_log

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type CkGormLogger struct {
	SlowThreshold time.Duration
}

func NewCkGormLogger() *CkGormLogger {
	return &CkGormLogger{
		SlowThreshold: 200 * time.Millisecond,
	}
}

var _ logger.Interface = (*CkGormLogger)(nil)

func (l *CkGormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return &CkGormLogger{}
}

func (l *CkGormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Infof(msg, data)
}
func (l *CkGormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *CkGormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *CkGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logFields := []logx.LogField{
		logx.Field("clickhouse sql", sql),
		logx.Field("time", MicrosecondsStr(elapsed)),
		logx.Field("rows", rows),
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.WithContext(ctx).Infow("Database ErrRecordNotFound", logFields...)
		} else {
			logFields = append(logFields, logx.Field("catch error", err))
			logx.WithContext(ctx).Errorw("Database Error", logFields...)
		}
	}
	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		logx.WithContext(ctx).Sloww("Database Slow Log", logFields...)
	}
	logx.WithContext(ctx).Infow("Database Query", logFields...)
}
