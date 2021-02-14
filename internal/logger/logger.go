package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

func NewLogger(lvlStr string, opts ...zap.Option) *zap.SugaredLogger {
	var lvl zapcore.Level
	if err := lvl.UnmarshalText([]byte(lvlStr)); err != nil {
		log.Fatal(err)
	}

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(lvl),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "timestamp",
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	l, err := config.Build(opts...)
	if err != nil {
		log.Fatal(err)
	}

	return l.Sugar()
}
