package logger

import (
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"seb7887/logsafe/masker"
	"strings"
	"time"
)

type zapLogger struct {
	log    *zap.Logger
	fields map[string]interface{}
	masker masker.Masker
}

func newZapLogger(cfg Config) (Logger, error) {
	// setup encoder
	c := zap.NewProductionEncoderConfig()
	c.MessageKey = cfg.Keys.MsgKey
	c.LevelKey = cfg.Keys.LevelKey
	c.TimeKey = cfg.Keys.TimeKey
	c.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.RFC3339))
	}
	enc := zapcore.NewJSONEncoder(c)

	// parse log level
	if cfg.Level == "" {
		cfg.Level = defaultLogLevel
	}
	l, err := zapcore.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		return nil, errors.New("invalid log level")
	}

	core := zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), l)
	log := zap.New(core)

	return &zapLogger{
		log:    log,
		fields: cfg.Fields,
		masker: masker.New(),
	}, nil
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.log.Debug(formatArgs(l.masker, args...), l.parseFields()...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.log.Info(formatArgs(l.masker, args...), l.parseFields()...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.log.Warn(formatArgs(l.masker, args...), l.parseFields()...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.log.Error(formatArgs(l.masker, args...), l.parseFields()...)
}

func (l *zapLogger) parseFields() []zap.Field {
	var fields []zap.Field
	for k, v := range l.fields {
		fields = append(fields, zap.Any(k, v))
	}
	return fields
}
