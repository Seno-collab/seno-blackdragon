package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log         *zap.Logger
	Sugar       *zap.SugaredLogger
	dailyWriter *DailyFileWriter
	level       zap.AtomicLevel
)

type LoggerConfig struct {
	Environment   string // "development" | "production"
	FilePath      string // e.g. "./logs/app.log" (leave empty = console only)
	Level         string // "debug" | "info" | "warn" | "error" (default: info)
	ConsolePretty bool   // true: pretty console output in dev
	DebugToFile   bool   // true: write debug logs to file as well
	// Lumberjack options (reasonable defaults if empty)
	MaxSizeMB   int  // default: 50 MB
	MaxBackups  int  // default: 7 files
	MaxAgeDays  int  // default: 14 days
	Compress    bool // default: true
	RotateDaily bool
}

// Init initializes the global logger according to the provided configuration.
// It can output to both console and file (if FilePath is set).
func Init(cfg LoggerConfig) error {
	// ----- Set log level -----
	level = zap.NewAtomicLevel()
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case "warn":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}

	// ----- Encoders -----
	// JSON encoder (used for file logging and production console)
	jsonEnc := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "msg",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		EncodeTime:     func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(time.RFC3339Nano)) },
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	})

	// Console encoder (pretty output in dev)
	var consoleEnc zapcore.Encoder
	if cfg.ConsolePretty || strings.ToLower(cfg.Environment) == "development" {
		ce := zap.NewDevelopmentEncoderConfig()
		ce.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { enc.AppendString(t.Format(time.RFC3339Nano)) }
		consoleEnc = zapcore.NewConsoleEncoder(ce)
	} else {
		consoleEnc = jsonEnc
	}

	// ----- Output targets -----
	consoleWS := zapcore.AddSync(os.Stdout)

	var cores []zapcore.Core
	// Console core
	cores = append(cores, zapcore.NewCore(consoleEnc, consoleWS, level))

	// File core (optional)
	if cfg.FilePath != "" {
		var fileWS zapcore.WriteSyncer
		if cfg.RotateDaily {
			dailyWriter = NewDailyFileWriter(cfg.FilePath)
			fileWS = zapcore.AddSync(dailyWriter)
		} else {
			maxSize := cfg.MaxSizeMB
			if maxSize <= 0 {
				maxSize = 50
			}
			maxBackups := cfg.MaxBackups
			if maxBackups <= 0 {
				maxBackups = 7
			}
			maxAge := cfg.MaxAgeDays
			if maxAge <= 0 {
				maxAge = 14
			}
			fileWS = zapcore.AddSync(&lumberjack.Logger{
				Filename:   cfg.FilePath, // Do not append date suffix; let lumberjack handle rotation
				MaxSize:    maxSize,      // in MB
				MaxBackups: maxBackups,
				MaxAge:     maxAge, // in days
				Compress:   true,
			})
		}

		fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			if cfg.DebugToFile {
				return lvl >= zap.DebugLevel && lvl >= level.Level()
			}
			// Only log Info+ to file unless DebugToFile is enabled
			return lvl >= zap.InfoLevel && lvl >= level.Level()
		})
		cores = append(cores, zapcore.NewCore(jsonEnc, fileWS, fileLevel))
	}

	// ----- Build logger -----
	core := zapcore.NewTee(cores...)
	l := zap.New(core,
		zap.AddCaller(),                   // add caller info (file:line)
		zap.AddStacktrace(zap.ErrorLevel), // include stacktrace for errors
		zap.AddCallerSkip(1),              // skip one caller level
		zap.ErrorOutput(consoleWS),        // output errors to console
		zap.WrapCore(func(c zapcore.Core) zapcore.Core { // sampling to reduce noise
			return zapcore.NewSamplerWithOptions(c, time.Second, 100, 100)
		}),
	)

	// Assign globals
	Log = l
	Sugar = l.Sugar()
	return nil
}

// SetLevel changes the log level at runtime (e.g. via an admin HTTP endpoint).
func SetLevel(lvl string) {
	switch strings.ToLower(lvl) {
	case "debug":
		level.SetLevel(zap.DebugLevel)
	case "warn":
		level.SetLevel(zap.WarnLevel)
	case "error":
		level.SetLevel(zap.ErrorLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}
}

// Close flushes buffered log entries and releases resources.
func Close() {
	if Log == nil {
		return
	}
	// Avoid "invalid argument" errors on Windows when syncing
	_ = Log.Sync()
}

type DailyFileWriter struct {
	basePath string
	curDate  string
	file     *os.File
	mu       sync.Mutex
}

func NewDailyFileWriter(basePath string) *DailyFileWriter {
	return &DailyFileWriter{basePath: basePath}
}

func (w *DailyFileWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if w.file == nil || today != w.curDate {
		if w.file != nil {
			w.file.Close()
		}
		filePath := fmt.Sprintf("%s-%s.log", w.basePath, today)
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.file = f
		w.curDate = today
	}
	return w.file.Write(p)
}
