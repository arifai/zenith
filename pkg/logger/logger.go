package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func InitLogger() {
	consoleSyncer := zapcore.AddSync(os.Stdout)

	infoFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "info.log",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	})

	errorFileSyncer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "error.log",
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
	})
	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, infoFileSyncer, infoLevel),
		zapcore.NewCore(jsonEncoder, errorFileSyncer, errorLevel),
		zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.DebugLevel),
	)

	Logger = zap.New(core, zap.AddCaller())
	defer func(Logger *zap.Logger) {
		_ = Logger.Sync()
	}(Logger)
}
