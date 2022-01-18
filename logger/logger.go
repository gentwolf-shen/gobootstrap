package logger

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gentwolf-shen/gobootstrap/helper/jsonhelper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Sugar *zap.SugaredLogger
)

func LoadDefault() {
	currentDir := filepath.Dir(os.Args[0])
	p := &LoggerStruct{}
	err := jsonhelper.FileToObj(currentDir+"/config/logger.json", p)
	if err != nil {
		panic(err)
	}

	p.FilePath = strings.Replace(p.FilePath, "${application.path}", currentDir, -1)
	if !strings.HasSuffix(p.FilePath, "/") {
		p.FilePath += "/"
	}

	p.LumberjackOption.Filename = p.FilePath + p.LumberjackOption.Filename
	zapcore.AddSync(&p.LumberjackOption)

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevel(zap.DebugLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(getWriter(p.FilePath, "info", &p.LumberjackOption)), getLevel(zapcore.InfoLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(getWriter(p.FilePath, "warn", &p.LumberjackOption)), getLevel(zapcore.WarnLevel)),
		zapcore.NewCore(encoder, zapcore.AddSync(getWriter(p.FilePath, "error", &p.LumberjackOption)), getLevel(zapcore.ErrorLevel)),
	)

	logger := zap.New(core, zap.AddCaller())
	defer logger.Sync()
	Sugar = logger.Sugar()
}

func getLevel(target zapcore.Level) zap.LevelEnablerFunc {
	return zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= target
	})
}

func getWriter(filePath string, level string, p *lumberjack.Logger) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filePath + level + ".log",
		MaxSize:    p.MaxSize,
		MaxAge:     p.MaxAge,
		MaxBackups: p.MaxBackups,
		LocalTime:  p.LocalTime,
		Compress:   p.Compress,
	}
}

type (
	LoggerStruct struct {
		FilePath         string            `json:"filepath"`
		LumberjackOption lumberjack.Logger `json:"lumberjackOption"`
	}
)
