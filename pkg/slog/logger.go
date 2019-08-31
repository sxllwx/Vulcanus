package slog

import (
	"fmt"
	"io"
	"log"
	"os"
)

// log level
const (
	info uint8 = iota
	debug
	err
	warn
	fatal
)

// color
const (
	red uint8 = iota + 91
	green
	yellow
	blue
	magenta
)

var paintBox = map[uint8]func(string) string{
	red:     func(s string) string { return colorful(red, s) },
	green:   func(s string) string { return colorful(green, s) },
	yellow:  func(s string) string { return colorful(yellow, s) },
	blue:    func(s string) string { return colorful(blue, s) },
	magenta: func(s string) string { return colorful(magenta, s) },
}

func colorful(color uint8, s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color, s)
}

var (
	levelToPrefix = map[uint8]string{
		info:  "[INFO] ",
		debug: "[DEBUG] ",
		err:   "[ERROR] ",
		warn:  "[WARN] ",
		fatal: "[FATAL] ",
	}

	levelToColor = map[uint8]func(string) string{
		info:  paintBox[blue],
		debug: paintBox[green],
		err:   paintBox[red],
		warn:  paintBox[yellow],
		fatal: paintBox[magenta],
	}

	dl = New(nil, WithColor())
)

// TODO add more opts,
// example split logfile
type options struct {
	color bool
}

// used to simple application
func WithColor() func(*options) {
	return func(o *options) {
		o.color = true
	}
}

type logger struct {
	*log.Logger
	o *options
}

type Logger interface {
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Err(...interface{})
	Errf(string, ...interface{})
}

func New(w io.Writer, opts ...func(*options)) *logger {

	var o = &options{}
	for _, f := range opts {
		f(o)
	}

	if w == nil {
		// not set, output to stdout
		w = os.Stdout
	}

	lr := log.New(w, "", log.LstdFlags|log.Llongfile|log.Ltime)
	return &logger{
		Logger: lr,
		o:      o,
	}
}

func (l *logger) output(calldepth int, msg string) {
	l.Logger.Output(calldepth, msg)
}

func (l *logger) assembleWithFormat(level uint8, format string, args ...interface{}) string {

	msg := levelToPrefix[level] + fmt.Sprintf(format, args...)
	if l.o.color {
		msg = levelToColor[level](msg)
	}
	return msg
}

func (l *logger) assemble(level uint8, args ...interface{}) string {

	msg := levelToPrefix[level] + fmt.Sprint(args...)
	if l.o.color {
		msg = levelToColor[level](msg)
	}
	return msg
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.output(3, l.assembleWithFormat(info, format, args...))
}

func (l *logger) Info(args ...interface{}) {
	l.output(3, l.assemble(info, args...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.output(3, l.assembleWithFormat(debug, format, args...))
}

func (l *logger) Debug(args ...interface{}) {
	l.output(3, l.assemble(debug, args...))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.output(3, l.assembleWithFormat(warn, format, args...))
}

func (l *logger) Warn(args ...interface{}) {
	l.output(3, l.assemble(warn, args...))
}

func (l *logger) Errf(format string, args ...interface{}) {
	l.output(3, l.assembleWithFormat(err, format, args...))
}

func (l *logger) Err(args ...interface{}) {
	l.output(3, l.assemble(err, args...))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.output(3, l.assembleWithFormat(fatal, format, args...))
	os.Exit(1)
}

func (l *logger) Fatal(args ...interface{}) {
	l.output(3, l.assemble(fatal, args...))
	os.Exit(1)
}

func Info(args ...interface{}) {
	dl.output(3, dl.assemble(info, args...))
}

func Infof(format string, args ...interface{}) {
	dl.output(3, dl.assembleWithFormat(info, format, args...))
}
func Debug(args ...interface{}) {
	dl.output(3, dl.assemble(debug, args...))
}

func Debugf(format string, args ...interface{}) {
	dl.output(3, dl.assembleWithFormat(debug, format, args...))
}

func Warn(args ...interface{}) {
	dl.output(3, dl.assemble(warn, args...))
}

func Warnf(format string, args ...interface{}) {
	dl.output(3, dl.assembleWithFormat(warn, format, args...))
}

func Err(args ...interface{}) {
	dl.output(3, dl.assemble(err, args...))
}

func Errf(format string, args ...interface{}) {
	dl.output(3, dl.assembleWithFormat(err, format, args...))
}

func Fatal(args ...interface{}) {
	dl.output(3, dl.assemble(fatal, args...))
	os.Exit(1)
}

func Fatalf(format string, args ...interface{}) {
	dl.output(3, dl.assembleWithFormat(fatal, format, args...))
	os.Exit(1)
}
