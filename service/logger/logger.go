package logger

import (
	"log"

	"github.com/cranemont/judge-manager/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Mode int
type Env string
type Level string

const (
	DEBUG Level = "Debug"
	INFO  Level = "Info"
	WARN  Level = "Warn"
	ERROR Level = "Error"
)

const (
	File Mode = 1 + iota
	Console
)

const (
	Production  Env = "production"
	Development Env = "development"
)

type Logger struct {
	zap *zap.Logger
}

func NewLogger(mode Mode, env Env) *Logger {
	var zapLogger *zap.Logger
	var cfg zap.Config
	var err error

	if env == Production {
		cfg = zap.NewProductionConfig()
		setMode(&cfg, mode, constants.LOG_PATH_PROD).
			EncoderConfig = zap.NewProductionEncoderConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapLogger, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	} else {
		cfg = zap.NewDevelopmentConfig()
		setMode(&cfg, mode, constants.LOG_PATH_DEV).
			EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		// cfg.Encoding = "json"
		zapLogger, err = cfg.Build(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}
	}
	return &Logger{zap: zapLogger}
}

func setMode(cfg *zap.Config, mode Mode, logPath string) *zap.Config {
	switch mode {
	case Console:
	case File:
		cfg.OutputPaths = []string{logPath}
	case File | Console:
		cfg.OutputPaths = append(cfg.OutputPaths, logPath)
	default:
		log.Fatalf("invalid logger mode: %d", mode)
	}
	return cfg
}

func (l *Logger) Log(level Level, msg string, fields ...zapcore.Field) {
	switch level {
	case DEBUG:
		l.Debug(msg, fields...)
	case INFO:
		l.Info(msg, fields...)
	case WARN:
		l.Warn(msg, fields...)
	case ERROR:
		l.Error(msg, fields...)
	}
}

func (l *Logger) Debug(msg string, fields ...zapcore.Field) {
	defer l.zap.Sync()
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zapcore.Field) {
	defer l.zap.Sync()
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zapcore.Field) {
	defer l.zap.Sync()
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zapcore.Field) {
	defer l.zap.Sync()
	l.zap.Error(msg, fields...)
}
