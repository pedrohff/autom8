package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(name string, debug bool) *zap.Logger {
	cfg := zap.NewProductionConfig()
	if !debug {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	logger, _ := cfg.Build()
	return logger.With(zap.String("app", name))
}

func main() {
	NewLogger("testtt", true).Debug("ae")
}
