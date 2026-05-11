package custom_log

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	SlowThreshold time.Duration
}

func NewGormLogger() *GormLogger {
	return &GormLogger{
		SlowThreshold: 200 * time.Millisecond,
	}
}

var _ logger.Interface = (*GormLogger)(nil)

func (l *GormLogger) LogMode(lev logger.LogLevel) logger.Interface {
	return &GormLogger{}
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Infof(msg, data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logx.WithContext(ctx).Errorf(msg, data)
}
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	logFields := []logx.LogField{
		logx.Field("sql", sql),
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

func MicrosecondsStr(d time.Duration) string {
	us := d.Microseconds()

	switch {
	case us < 1000:
		return fmt.Sprintf("%dμs", us)
	case us < 1_000_000:
		return fmt.Sprintf("%dms", us/1000)
	default:
		return fmt.Sprintf("%.2fs", float64(us)/1_000_000)
	}
}
