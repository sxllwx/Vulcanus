package slog

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {

	Info("hello", "scott")
	Infof("hello %s", "scott")

	Debug("hello", "scott")
	Debugf("hello %s", "scott")

	Warn("hello", "scott")
	Warnf("hello %s", "scott")

	Err("hello", "scott")
	Errf("hello %s", "scott")

	//	Fatal("hello", "scott")
	//	Fatalf("hello %s", "scott")
}

func TestNew(t *testing.T) {

	f, err := os.OpenFile("test.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}
	l := New(f, WithColor())
	l.Info("hello", "scott")
	l.Infof("hello %s", "scott")

	l.Debug("hello", "scott")
	l.Debugf("hello %s", "scott")

	l.Warn("hello", "scott")
	l.Warnf("hello %s", "scott")

	l.Err("hello", "scott")
	l.Errf("hello %s", "scott")

	if err := os.Remove("test.log"); err != nil {
		t.Fatal(err)
	}

}
