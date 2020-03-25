package log

import (
	"testing"

	"go.uber.org/zap"
)

func TestLogger(t *testing.T) {

	Info("hello", zap.String("id", "1234"))
	Infof("hello %s", "scott")

	Debug("hello", zap.String("id", "1234"))
	Debugf("hello %s", "scott")

	Error("hello i am a err, should trace me")
	Pa("hello i am a err, should not  trace me")
}
