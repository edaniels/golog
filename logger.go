package golog

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
)

type Logger = *zap.SugaredLogger

var Global Logger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Global = logger.Sugar()
}

// NewTestLogger directs logs to the go test logger.
func NewTestLogger(t *testing.T) Logger {
	logger, _ := NewObservedTestLogger(t)
	return logger
}

// NewObservedTestLogger is like NewTestLogger but also saves logs to an in memory observer.
func NewObservedTestLogger(t *testing.T) (Logger, *observer.ObservedLogs) {
	logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.AddCaller()))
	observerCore, observedLogs := observer.New(zap.LevelEnablerFunc(zapcore.DebugLevel.Enabled))
	logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(c, observerCore)
	}))
	return logger.Sugar(), observedLogs
}
