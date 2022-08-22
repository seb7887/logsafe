package logger

import (
	"errors"
)

// Logger provides fast, structured and levered logging.
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// Keys logging keys settings
type Keys struct {
	MsgKey   string
	LevelKey string
	TimeKey  string
}

// Config logger settings
type Config struct {
	Fields map[string]interface{}
	Keys   Keys
	Level  string
	Type   int
}

const (
	// LogrusLogger logrus type Logger
	LogrusLogger = iota
	// ZapLogger zap type Logger
	ZapLogger
	defaultLogLevel = "DEBUG"
)

var (
	// ErrInvalidConfig error returned when an invalid config is provided for the Logger instance.
	ErrInvalidConfig = errors.New("invalid logger configuration")
)

// New returns a Logger instance
func New(cfg ...Config) (Logger, error) {
	if len(cfg) > 1 {
		return nil, ErrInvalidConfig
	}

	defaultCfg := Config{
		Keys: Keys{
			LevelKey: "level",
			MsgKey:   "msg",
			TimeKey:  "ts",
		},
		Level: defaultLogLevel,
		Type:  LogrusLogger,
	}

	if len(cfg) == 0 {
		cfg = []Config{defaultCfg}
	}

	switch cfg[0].Type {
	case LogrusLogger:
		return newLogrusLogger(cfg[0])
	case ZapLogger:
		return newZapLogger(cfg[0])
	default:
		return nil, ErrInvalidConfig
	}
}
