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
	logger, err := defaultDevelopmentConfig.Build()
	if err != nil {
		panic(err)
	}
	Global = logger.Sugar()
}

// from https://github.com/uber-go/zap/blob/2314926ec34c23ee21f3dd4399438469668f8097/config.go#L98
// but disable stacktraces.
var defaultProductionConfig = zap.Config{
	Level:       NewAtomicLevelAt(DebugLevel),
	Development: true,
	Encoding:    "json",
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},
	DisableStacktrace: true,
	OutputPaths:       []string{"stdout"},
	ErrorOutputPaths:  []string{"stderr"},
}

// from https://github.com/uber-go/zap/blob/2314926ec34c23ee21f3dd4399438469668f8097/config.go#L135
// but disable stacktraces, use same keys as prod, and color levels.
var defaultDevelopmentConfig = zap.Config{
	Level:    NewAtomicLevelAt(DebugLevel),
	Encoding: "console",
	EncoderConfig: zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	},
	DisableStacktrace: true,
	OutputPaths:       []string{"stdout"},
	ErrorOutputPaths:  []string{"stderr"},
}

// NewLogger returns a new logger using the default production configuration.
func NewLogger(name string) Logger {
	logger, err := defaultProductionConfig.Build()
	if err != nil {
		Global.Fatal(err)
	}
	return logger.Sugar().Named(name)
}

// NewDevelopmentLogger returns a new logger using the default development configuration.
func NewDevelopmentLogger(name string) Logger {
	logger, err := defaultDevelopmentConfig.Build()
	if err != nil {
		Global.Fatal(err)
	}
	return logger.Sugar().Named(name)
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
