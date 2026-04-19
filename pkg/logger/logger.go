package logger

import (
	"go-learning/pkg/settings"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config settings.LoggerSetting) *LoggerZap {
	// TODO: Init logger from file or env
	encoder := getEncoderLog()
	hook := lumberjack.Logger{
		Filename:   config.Path,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		getLogLevel("debug"))

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))}
}

func getLogLevel(level string) zapcore.Level {
	// debug -> info -> warn -> error -> fatal -> panic
	var res zapcore.Level

	switch level {
	case "debug":
		res = zapcore.DebugLevel
	case "info":
		res = zapcore.InfoLevel
	case "warn":
		res = zapcore.WarnLevel
	case "error":
		res = zapcore.ErrorLevel
	default:
		res = zapcore.InfoLevel
	}

	return res
}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	// 1776416764.2753186 -> 2026-04-17T16:06:04.275+0700
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// ts -> time
	encodeConfig.TimeKey = "time"
	// info -> INFO
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// cli/main.log.go:21
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}
