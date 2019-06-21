package slog

import (
	"fmt"
	"log"
	"os"
)

const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红
	info          = "[INFO] "
	debug         = "[DEBUG] "
	err           = "[ERRO] "
	warn          = "[WARN] "
	study         = "[STUDY] "
)

type Study interface {
	Study(args ...interface{})
	Studyf(format string, args ...interface{})
}

type Logger struct {
	l *log.Logger
}

func (l *Logger) Info(args ...interface{}) {
	l.Infof("%v", args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Warningf("%v", args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Errorf("%v", args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.Debugf("%v", args...)
}

func (l *Logger) Study(args ...interface{}) {
	l.Studyf("%v", args...)
}

func (l *Logger) Warnf(fmt string, args ...interface{}) {
	l.Warningf(fmt, args...)
}

func New() *Logger {
	lr := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Ltime)
	return &Logger{l: lr}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.Output(4, blue(info+fmt.Sprintf(format, args...)))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.Output(4, green(debug+fmt.Sprintf(format, args...)))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.l.Output(4, red(err+fmt.Sprintf(format, args...)))
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.l.Output(4, yellow(warn+fmt.Sprintf(format, args...)))
}

func (l *Logger) Studyf(format string, args ...interface{}) {
	l.l.Output(4, magenta(study+fmt.Sprintf(format, args...)))
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}
