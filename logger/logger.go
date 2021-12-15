package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var (
	Sugar *zap.SugaredLogger
)

func SetLoggerLevel(level zapcore.Level) {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	core := zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), os.Stdout, level)
	logger := zap.New(core, zap.AddCaller()).WithOptions()
	defer logger.Sync()
	Sugar = logger.Sugar()
}
