package util

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func zapInit() {
	var logger *zap.Logger

	c := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	logPath := "./zap.log"
	MaxSize := 100
	MaxAge := 1
	fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename: logPath, // 日誌文件存放目錄，如果文件夾不存在會自動創建
		MaxSize:  MaxSize, // 文件大小限制,單位MB
		MaxAge:   MaxAge,  // 日誌文件保留天數
		Compress: false,   // 是否壓縮處理
	})

	level := zapcore.InfoLevel

	core := zapcore.NewTee(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(c),
			zapcore.NewMultiWriteSyncer(
				zapcore.AddSync(os.Stdout),
				fileWriteSyncer),
			level),
	)

	// caller 顯示文件名、行號和zap調用者的函數名
	if level == zapcore.DebugLevel {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}

	logger.Sugar()
}
