package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level            string   `mapstructure:"level"`
	OutputPaths      []string `mapstructure:"output_paths"`
	ErrorOutputPaths []string `mapstructure:"error_output_paths"`
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type logger struct {
	log *zap.Logger
}

func New(conf *Config, opts ...zap.Option) (Logger, error) {
	level, err := zap.ParseAtomicLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	l, err := zap.Config{
		Level:       level,
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "ts",
			LevelKey:    "level",
			MessageKey:  "message",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.EpochTimeEncoder,
		},
		OutputPaths:      conf.OutputPaths,
		ErrorOutputPaths: conf.ErrorOutputPaths,
	}.Build(opts...)
	if err != nil {
		return nil, err
	}
	return &logger{
		log: l,
	}, nil
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

func (l logger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l logger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

type MockLogger struct {
	logs map[string][]bodyMessage
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		logs: make(map[string][]bodyMessage),
	}
}

type bodyMessage struct {
	msg    string
	fields map[string]string
}

func (t *MockLogger) Debug(msg string, fields ...zap.Field) {
	t.addMessage("debug", msg, fields...)
}

func (t *MockLogger) Info(msg string, fields ...zap.Field) {
	t.addMessage("info", msg, fields...)
}

func (t *MockLogger) Error(msg string, fields ...zap.Field) {
	t.addMessage("error", msg, fields...)
}

func (t *MockLogger) Fatal(msg string, fields ...zap.Field) {
	t.addMessage("fatal", msg, fields...)
}

func (t *MockLogger) addMessage(level, msg string, fields ...zap.Field) {
	l := t.logs[level]
	if l == nil {
		logFields := make(map[string]string)
		for _, field := range fields {
			v := fmt.Sprintf("%d", field.Integer)
			if field.String != "" {
				v = field.String
			}
			logFields[field.Key] = v
		}
		l = []bodyMessage{
			{msg: msg, fields: logFields},
		}
	}
	t.logs[level] = l
}

func (t *MockLogger) HasMessage(level, msg, key, value string) bool {
	l := t.logs[level]
	if l == nil {
		return false
	}
	for _, v := range l {
		if v.msg == msg {
			if f, ok := v.fields[key]; ok && f == value {
				return true
			}
		}
	}
	return false
}
