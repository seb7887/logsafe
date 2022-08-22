package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"reflect"
	"seb7887/logsafe/masker"
	"testing"
)

var (
	cfg = Config{
		Fields: nil,
		Keys: Keys{
			TimeKey:  "ts",
			MsgKey:   "msg",
			LevelKey: "level",
		},
		Level: "",
		Type:  ZapLogger,
	}
)

func TestZapLogger_New(t *testing.T) {
	l, err := newZapLogger(cfg)

	require.NoError(t, err, "should not return error")
	assert.Equal(t, "*logger.zapLogger", reflect.ValueOf(l).Type().String())
}

func TestZapLogger_Debug(t *testing.T) {
	obsCore, obsLogs := observer.New(zap.DebugLevel)
	log := zap.New(obsCore)

	l := zapLogger{
		log:    log,
		fields: nil,
		masker: masker.New(),
	}

	l.Debug("test")

	assert.Equal(t, 1, obsLogs.Len())
	firstLog := obsLogs.All()[0]
	assert.Equal(t, "test", firstLog.Message)
}

func TestZapLogger_Info(t *testing.T) {
	obsCore, obsLogs := observer.New(zap.DebugLevel)
	log := zap.New(obsCore)

	l := zapLogger{
		log:    log,
		fields: nil,
		masker: masker.New(),
	}

	l.Info("test")

	assert.Equal(t, 1, obsLogs.Len())
	firstLog := obsLogs.All()[0]
	assert.Equal(t, "test", firstLog.Message)
}

func TestZapLogger_Warn(t *testing.T) {
	obsCore, obsLogs := observer.New(zap.DebugLevel)
	log := zap.New(obsCore)

	l := zapLogger{
		log:    log,
		fields: nil,
		masker: masker.New(),
	}

	l.Warn("test")

	assert.Equal(t, 1, obsLogs.Len())
	firstLog := obsLogs.All()[0]
	assert.Equal(t, "test", firstLog.Message)
}

func TestZapLogger_Error(t *testing.T) {
	obsCore, obsLogs := observer.New(zap.DebugLevel)
	log := zap.New(obsCore)

	l := zapLogger{
		log:    log,
		fields: nil,
		masker: masker.New(),
	}

	l.Error("test")

	assert.Equal(t, 1, obsLogs.Len())
	firstLog := obsLogs.All()[0]
	assert.Equal(t, "test", firstLog.Message)
}
