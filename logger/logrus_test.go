package logger

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"reflect"
	"seb7887/logsafe/masker"
	"strings"
	"testing"
)

var (
	log = logrus.New()
)

func TestLogrusLogger_New(t *testing.T) {
	cfg := Config{}
	l, err := newLogrusLogger(cfg)

	require.NoError(t, err, "should not return error")
	assert.Equal(t, "*logger.logrusLogger", reflect.ValueOf(l).Type().String())
}

func TestLogrusLogger_Debug(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = logrusLogger{
			log:    log,
			fields: nil,
			masker: masker.New(),
		}
	)

	log.SetOutput(&buf)
	log.SetLevel(logrus.DebugLevel)
	l.Debug("test")
	assert.True(t, strings.Contains(buf.String(), "test"))
	log.SetOutput(os.Stderr)
}

func TestLogrusLogger_Info(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = logrusLogger{
			log:    log,
			fields: nil,
			masker: masker.New(),
		}
	)

	log.SetOutput(&buf)
	log.SetLevel(logrus.InfoLevel)
	l.Info("test")
	assert.True(t, strings.Contains(buf.String(), "test"))
	log.SetOutput(os.Stderr)
}

func TestLogrusLogger_Warn(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = logrusLogger{
			log:    log,
			fields: nil,
			masker: masker.New(),
		}
	)

	log.SetOutput(&buf)
	log.SetLevel(logrus.InfoLevel)
	l.Warn("test")
	assert.True(t, strings.Contains(buf.String(), "test"))
	log.SetOutput(os.Stderr)
}

func TestLogrusLogger_Error(t *testing.T) {
	var (
		buf bytes.Buffer
		l   = logrusLogger{
			log:    log,
			fields: nil,
			masker: masker.New(),
		}
	)

	log.SetOutput(&buf)
	log.SetLevel(logrus.InfoLevel)
	l.Error("test")
	assert.True(t, strings.Contains(buf.String(), "test"))
	log.SetOutput(os.Stderr)
}
