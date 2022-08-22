package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"seb7887/logsafe/masker"
	"time"
)

type logrusLogger struct {
	log    *logrus.Logger
	fields map[string]interface{}
	masker masker.Masker
}

func newLogrusLogger(cfg Config) (Logger, error) {
	l := logrus.New()
	l.SetOutput(os.Stdout)

	if cfg.Level == "" {
		cfg.Level = defaultLogLevel
	}
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}
	l.SetLevel(level)

	formatter := &logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyLevel: cfg.Keys.LevelKey,
			logrus.FieldKeyMsg:   cfg.Keys.MsgKey,
			logrus.FieldKeyTime:  cfg.Keys.TimeKey,
		},
		TimestampFormat: time.RFC3339,
	}

	l.SetFormatter(formatter)

	return &logrusLogger{
		log:    l,
		fields: cfg.Fields,
		masker: masker.New(),
	}, nil
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.log.WithFields(l.fields).Debug(formatArgs(l.masker, args...))
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.log.WithFields(l.fields).Info(formatArgs(l.masker, args...))
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.log.WithFields(l.fields).Warn(formatArgs(l.masker, args...))
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.log.WithFields(l.fields).Error(formatArgs(l.masker, args...))
}
