package log

import (
	"io"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// global sys var
var (
	GlobalLogger *defaultLogger
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
	Panicf func(template string, args ...interface{})
	Fatalf func(template string, args ...interface{})
)

func init() {

	GlobalLogger = newLogger()
	GlobalLogger.installGoStyleLogFunc()
	GlobalLogger.installZapLogFunc()

	// the sync method
	Sync = GlobalLogger.Sync
}

type Config struct {
	EncoderCfg zapcore.EncoderConfig
	Encoder    zapcore.Encoder
	Writers    zapcore.WriteSyncer
	Level      zap.AtomicLevel

	Core zapcore.Core
}

type defaultLogger struct {
	*zap.Logger
	productEnvConfig Config
	developEnvConfig Config
}

func newLogger() *defaultLogger {

	out := &defaultLogger{}

	// default set to discard
	out.applyProduct(zap.ErrorLevel, ioutil.Discard)
	// default set to os.Stdout
	out.applyDevelop(zap.DebugLevel, os.Stdout)

	out.Logger = zap.New(
		zapcore.NewTee(out.productEnvConfig.Core, out.developEnvConfig.Core),
		zap.AddCaller(),
		//zap.AddStacktrace(zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		//	return l >= zap.WarnLevel
		//})),
	)

	return out
}

func (l *defaultLogger) applyProduct(level zapcore.Level, writer io.Writer) {

	// direct use zap std config
	l.productEnvConfig.EncoderCfg = zap.NewProductionEncoderConfig()
	l.productEnvConfig.Encoder = zapcore.NewJSONEncoder(l.productEnvConfig.EncoderCfg)

	l.productEnvConfig.Level = zap.NewAtomicLevel()
	l.productEnvConfig.Level.SetLevel(level)

	l.productEnvConfig.Writers = zapcore.AddSync(writer)

	// set full core
	l.productEnvConfig.Core = zapcore.NewCore(
		l.productEnvConfig.Encoder,
		l.productEnvConfig.Writers,
		l.productEnvConfig.Level,
	)
}

func (l *defaultLogger) applyDevelop(level zapcore.Level, writer io.Writer) {

	l.developEnvConfig.EncoderCfg = zap.NewDevelopmentEncoderConfig()

	// color for the develop
	l.developEnvConfig.EncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// short name will more convenient
	// l.developEnvConfig.EncoderCfg.EncodeCaller = zapcore.FullCallerEncoder

	l.developEnvConfig.Encoder = zapcore.NewConsoleEncoder(l.developEnvConfig.EncoderCfg)

	l.developEnvConfig.Level = zap.NewAtomicLevel()
	l.developEnvConfig.Level.SetLevel(level)

	l.developEnvConfig.Writers = zapcore.AddSync(writer)

	// set full core
	l.developEnvConfig.Core = zapcore.NewCore(
		l.developEnvConfig.Encoder,
		l.developEnvConfig.Writers,
		l.developEnvConfig.Level,
	)
}

func (l *defaultLogger) DevelopLevel() zap.AtomicLevel {
	return l.developEnvConfig.Level
}

func (l *defaultLogger) ProductLevel() zap.AtomicLevel {
	return l.productEnvConfig.Level
}

func (l *defaultLogger) installZapLogFunc() {

	Info = l.Info
	Debug = l.Debug
	Warn = l.Warn
	Error = l.Error
	Fatal = l.Fatal
}

func (l *defaultLogger) installGoStyleLogFunc() {

	Infof = l.Sugar().Infof
	Warnf = l.Sugar().Warnf
	Debugf = l.Sugar().Debugf
	Errorf = l.Sugar().Errorf
	Fatalf = l.Sugar().Fatalf
	Panicf = l.Sugar().Panicf
}
