package logger

import (
	"leaderboard/config"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger -
func NewZapLogger(conf config.Config) *zap.Logger {
	var zapLevel zapcore.Level

	switch strings.ToLower(conf.Logger) {
	case "info":
		zapLevel = zapcore.InfoLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "warning":
		zapLevel = zapcore.WarnLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	case "debug":
		zapLevel = zapcore.DebugLevel
	default:
		zapLevel = zapcore.DebugLevel
	}

	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.Lock(os.Stderr),
		zapLevel,
	)

	return zap.New(core, zap.AddCallerSkip(2), zap.AddStacktrace(zapcore.PanicLevel))
}
