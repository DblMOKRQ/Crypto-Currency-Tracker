package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

// NewLogger создает и настраивает новый экземпляр логгера
func NewLogger(logLevel string) (*zap.Logger, error) {

	var level zapcore.Level
	var encoding string
	var encodeLevel zapcore.LevelEncoder
	switch logLevel {
	case "debug":
		level = zapcore.DebugLevel
		encoding = "console"
		encodeLevel = zapcore.CapitalColorLevelEncoder
	case "info":
		level = zapcore.InfoLevel
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	case "warn":
		level = zapcore.WarnLevel
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	case "error":
		level = zapcore.ErrorLevel
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	default:
		level = zapcore.InfoLevel // По умолчанию уровень Info
		encoding = "json"
		encodeLevel = zapcore.LowercaseLevelEncoder
	}

	// Настройка конфигурации логгера
	config := zap.Config{
		Encoding:         encoding,
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     customTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	// Создание логгера
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("error building zap logger: %w", err)
	}

	logger.Info("Logger initialized",
		zap.String("level", level.String()),
		zap.String("encoding", config.Encoding))

	return logger, nil
}

// Кастомный формат времени
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02-01-2006 15:04:05"))
}
