package logger

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestLogger_NewWithConfigLogrus(t *testing.T) {
	cfg := Config{}
	l, err := New(cfg)

	require.NoError(t, err, "should not return error")
	assert.Equal(t, "*logger.logrusLogger", reflect.ValueOf(l).Type().String())
}

func TestLogger_NewWithConfigZap(t *testing.T) {
	cfg := Config{
		Type: ZapLogger,
	}
	l, err := New(cfg)

	require.NoError(t, err, "should not return error")
	assert.Equal(t, "*logger.zapLogger", reflect.ValueOf(l).Type().String())
}

func TestLogger_NewWithNoConfig(t *testing.T) {
	l, err := New()

	require.NoError(t, err, "should not return error")
	assert.Equal(t, "*logger.logrusLogger", reflect.ValueOf(l).Type().String())
}

func TestLogger_NewWithInvalidConfig(t *testing.T) {
	cfg1 := Config{}
	cfg2 := Config{}
	_, err := New(cfg1, cfg2)

	require.Error(t, err, "should return error")
	assert.Equal(t, ErrInvalidConfig, err)
}
