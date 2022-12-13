package logger

import (
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
