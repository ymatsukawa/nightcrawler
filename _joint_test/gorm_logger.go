package main

// ref. https://gorm.io/docs/logger.html
import (
	"context"
	"log/slog"
	"time"

	glogger "gorm.io/gorm/logger"
)

type slogGormLogger struct {
	logger *slog.Logger
	level  glogger.LogLevel
}

func newSlogGormLogger(l *slog.Logger) *slogGormLogger {
	return &slogGormLogger{logger: l, level: glogger.Info}
}

func (l *slogGormLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	nl := *l
	nl.level = level
	return &nl
}

func (l *slogGormLogger) Info(ctx context.Context, msg string, args ...interface{})  {}
func (l *slogGormLogger) Warn(ctx context.Context, msg string, args ...interface{})  {}
func (l *slogGormLogger) Error(ctx context.Context, msg string, args ...interface{}) {}

func (l *slogGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= glogger.Silent {
		return
	}

	sql, rows := fc()
	if sql == "" {
		return
	}

	l.logger.LogAttrs(ctx, slog.LevelInfo, sql,
		slog.Int64("rows", rows),
		slog.Duration("elapsed", time.Since(begin)),
	)
}
