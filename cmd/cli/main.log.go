package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// sugar := zap.NewExample().Sugar()
	// defer sugar.Sync()
	// sugar.Infof("Hello, name: %s, age: %d", "Phuc", 25)

	// logger := zap.NewExample()
	// defer logger.Sync()
	// logger.Info("Hello", zap.String("name", "Phuc"), zap.Int("age", 25))

	// logger := zap.NewExample()
	// logger.Info("Hello")

	// logger, _ = zap.NewDevelopment()
	// logger.Info("Hello Development")

	// logger, _ = zap.NewProduction()
	// logger.Info("Hello Production")

	encoder := getEncoderLog()
	sync := getWriterSync()
	core := zapcore.NewCore(encoder, sync, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())

	logger.Info("Info log", zap.Int("line", 1))
	logger.Warn("Warn log", zap.Int("line", 2))
	logger.Error("Error log", zap.Int("line", 3))
}

// format log
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

func getWriterSync() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	syncFile := zapcore.AddSync(file)
	syncConsole := zapcore.AddSync(os.Stderr)
	return zap.CombineWriteSyncers(syncFile, syncConsole)
}
