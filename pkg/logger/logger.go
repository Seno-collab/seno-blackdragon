package logger

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

type LoggerConfig struct {
	Environment string
	FilePath    string
}

func getLogWriterByDate(basePath string) zapcore.WriteSyncer {
	today := time.Now().Format("2006-01-02")
	fullPath := fmt.Sprintf("%s-%s.log", basePath, today)
	lumberjackLogger := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
	return zapcore.AddSync(lumberjackLogger)
}

func Init(cfg LoggerConfig) {
	var zapCfg zap.Config
	if cfg.Environment == "development" {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.FilePath != "" {
		writer := getLogWriterByDate(cfg.FilePath)
		core := zapcore.NewCore(zapcore.NewJSONEncoder(zapCfg.EncoderConfig), writer, zap.NewAtomicLevelAt(zap.InfoLevel))
		Log = zap.New(core)
	} else {
		logger, _ := zapCfg.Build()
		Log = logger
	}
	Sugar = Log.Sugar()
}

func Close() {
	_ = Log.Sync()
}
