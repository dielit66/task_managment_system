package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(level int8) (ILogger, error) {
	var zapLevel zapcore.Level

	switch level {
	case -1:
		zapLevel = zap.DebugLevel
	case 0:
		zapLevel = zap.InfoLevel
	case 1:
		zapLevel = zap.WarnLevel
	case 2:
		zapLevel = zap.ErrorLevel
	}

	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeTime:  zapcore.ISO8601TimeEncoder,
			EncodeLevel: zapcore.CapitalLevelEncoder,
		},
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{logger: zapLogger}, nil
}

func (l *ZapLogger) Debug(msg string, fields ...interface{}) {
	l.logger.Debug(msg, convertFields(fields...)...)
}

func (l *ZapLogger) Info(msg string, fields ...interface{}) {
	l.logger.Info(msg, convertFields(fields...)...)
}

func (l *ZapLogger) Warn(msg string, fields ...interface{}) {
	l.logger.Warn(msg, convertFields(fields...)...)
}

func (l *ZapLogger) Error(msg string, fields ...interface{}) {
	l.logger.Error(msg, convertFields(fields...)...)
}

func (l *ZapLogger) Fatal(msg string, fields ...interface{}) {
	l.logger.Fatal(msg, convertFields(fields...)...)
}

func convertFields(fields ...interface{}) []zap.Field {
	var zapFields []zap.Field
	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		zapFields = append(zapFields, zap.Any(key, fields[i+1]))
	}
	return zapFields
}
