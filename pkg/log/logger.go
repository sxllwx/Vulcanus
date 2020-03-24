package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global sys var
var (
	GlobalLogger *zap.Logger
	Level        zap.AtomicLevel
	LogDir       string
)

var (
	// after the app close, we should sync the message to io.Writer
	Sync func() error
)

var (
	Info  func(msg string, fields ...zap.Field)
	Debug func(msg string, fields ...zap.Field)
	Error func(msg string, fields ...zap.Field)
	Warn  func(msg string, fields ...zap.Field)
	Fatal func(msg string, fields ...zap.Field)

	// work like stand pkg sys.Printf
	Infof  func(template string, args ...interface{})
	Warnf  func(template string, args ...interface{})
	Debugf func(template string, args ...interface{})
	Errorf func(template string, args ...interface{})
	Fataf  func(template string, args ...interface{})
)

func newStdOutLogCore() zapcore.Core {

	// std use the console encoder
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		Level,
	)
}

func newLogFileLogCore() zapcore.Core {

	// Set the fd
	if len(LogDir) == 0 {
		return zapcore.NewNopCore()
	}
	logFileFD, err := os.OpenFile(LogDir, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	// sys-file the console encoder
	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(logFileFD),
		Level,
	)
}

func init() {

	Level = zap.NewAtomicLevel()

	trace := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zap.WarnLevel
	})

	GlobalLogger = zap.New(zapcore.NewTee(
		newStdOutLogCore(),
		newLogFileLogCore(),
	),
		zap.AddCaller(),
		zap.AddStacktrace(trace),
	)

	Info = GlobalLogger.Info
	Debug = GlobalLogger.Debug
	Warn = GlobalLogger.Warn
	Error = GlobalLogger.Error
	Fatal = GlobalLogger.Fatal

	// work like stand pkg sys.Printf
	Infof = GlobalLogger.Sugar().Infof
	Warnf = GlobalLogger.Sugar().Warnf
	Debugf = GlobalLogger.Sugar().Debugf
	Errorf = GlobalLogger.Sugar().Errorf
	Fataf = GlobalLogger.Sugar().Fatalf

	// the sync method
	Sync = GlobalLogger.Sync
}
