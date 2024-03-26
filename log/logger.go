package log

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/JamesClonk/item-monitor/config"
)

var (
	logger *slog.Logger
)

func init() {
	logger = newLogger(os.Stdout)
}

func newLogger(writer io.Writer) *slog.Logger {
	logLevel := slog.LevelDebug
	switch config.Get().LogLevel {
	case "info":
		logLevel = slog.LevelInfo
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(writer, opts))

	return logger
}

func Info(format string, args ...interface{}) {
	logger.Info(format, args...)
}
func Infof(format string, args ...interface{}) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	logger.Warn(format, args...)
}
func Warnf(format string, args ...interface{}) {
	logger.Warn(fmt.Sprintf(format, args...))
}

func Debug(format string, args ...interface{}) {
	logger.Debug(format, args...)
}
func Debugf(format string, args ...interface{}) {
	logger.Debug(fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	logger.Error(format, args...)
}
func Errorf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Fatal(format string, args ...interface{}) {
	logger.Error(format, args...)
	os.Exit(1)
}
func Fatalf(format string, args ...interface{}) {
	logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}
