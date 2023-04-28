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

func LoadFromStr(str string) {
	p := &LoggerStruct{}
	if err := jsonhelper.ToObj([]byte(str), p); err != nil {
		panic(err)
	}

	initLogger(p)
}

func LoadFromFile(filename string) {
	p := &LoggerStruct{}
	if err := jsonhelper.FileToObj(filename, p); err != nil {
		panic(err)
	}

	initLogger(p)
}

func LoadDefault() {
	LoadFromFile(filepath.Dir(os.Args[0]) + "/config/logger.json")
}

func initLogger(p *LoggerStruct) {
	p.FilePath = strings.Replace(p.FilePath, "${application.path}", filepath.Dir(os.Args[0]), -1)
	if !strings.HasSuffix(p.FilePath, "/") {
		p.FilePath += "/"
	}

	p.LumberjackOption.Filename = p.FilePath + p.LumberjackOption.Filename
	zapcore.AddSync(&p.LumberjackOption)

	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)

	index := 1
	cores := make([]zapcore.Core, 4)
	cores[0] = zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), getLevel(p.Level))
	levels := []zapcore.Level{zapcore.InfoLevel, zapcore.WarnLevel, zapcore.ErrorLevel}
	for _, level := range levels {
		if level >= p.Level {
			cores[index] = zapcore.NewCore(encoder, zapcore.AddSync(getWriter(p.FilePath, level.String(), &p.LumberjackOption)), getLevel(level))
			index++
		}
	}

	core := zapcore.NewTee(cores[0:index]...)
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
		Level            zapcore.Level     `json:"level"`
		LumberjackOption lumberjack.Logger `json:"lumberjackOption"`
	}
)
