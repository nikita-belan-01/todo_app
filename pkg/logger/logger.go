package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKey string

var ErrInvalidLogLevel = errors.New("invalid log level")

const (
	LoggerKeyName loggerKey = "log"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

func ValidateLogLevel(logLevel string) error {
	if logLevel != zap.DebugLevel.String() &&
		logLevel != zap.InfoLevel.String() &&
		logLevel != zap.WarnLevel.String() &&
		logLevel != zap.ErrorLevel.String() &&
		logLevel != zap.DPanicLevel.String() &&
		logLevel != zap.PanicLevel.String() &&
		logLevel != zap.FatalLevel.String() {
		return ErrInvalidLogLevel
	}

	return nil
}

func New(logDir, logLevel string) (*Logger, error) {
	zapLevel := zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("create log dir: %w", err)
	}

	nowStr := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		logDir,
		fmt.Sprintf("%s.log", nowStr))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("create log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLevel),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLevel))

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		zapLogger,
		logFile}, nil
}

func (l Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		l.Logger.With(fields...),
		l.file}
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(LoggerKeyName).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

func (l *Logger) Close() error {
	var fileSyncErr, fileCloseErr error

	_ = l.Sync()

	if err := l.file.Sync(); err != nil {
		fileSyncErr = fmt.Errorf("sync logger file: %w", err)
	}

	if err := l.file.Close(); err != nil {
		fileCloseErr = fmt.Errorf("close logger file: %w", err)
	}

	return errors.Join(fileSyncErr, fileCloseErr)
}
