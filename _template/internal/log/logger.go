package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewAtLevel(levelStr string) (*zap.Logger, error) {
	logLevel := zapcore.InfoLevel
	if levelStr != "" {
		var err error
		logLevel, err = zapcore.ParseLevel(levelStr)
		if err != nil {
			return nil, err
		}
	}

	logConf := zap.NewProductionConfig()
	logConf.Level = zap.NewAtomicLevelAt(logLevel)

	logger, err := logConf.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}
