package slog

import (
	"fmt"
	"log"
	"os"
)

const (
	INFO    = "INFO "
	DEBUG   = "DEBUG "
	WARNING = "WARNING "
	ERROR   = "ERROR "
)

type Logger struct {
	module string
	l      *log.Logger
}

func New(module string) *Logger {
	lr := log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	return &Logger{l: lr, module: module}
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.Output(2, l.modulePrefix()+INFO+fmt.Sprintf(format, args...))
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.Output(2, l.modulePrefix()+DEBUG+fmt.Sprintf(format, args...))
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.l.Output(2, l.modulePrefix()+ERROR+fmt.Sprintf(format, args...))
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.l.Output(2, l.modulePrefix()+WARNING+fmt.Sprintf(format, args...))
}

func (l *Logger) modulePrefix() string {
	return fmt.Sprintf(" [ %s ] ", l.module)
}
