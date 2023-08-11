package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitZap(levelStr string) (*zap.Logger, error) {
	var level zapcore.Level
	err := level.Set(levelStr)
	if err != nil {
		return nil, err
	}

	loggerConfig := zap.NewProductionConfig()
	if err != nil {
		return nil, err
	}

	loggerConfig.Level.SetLevel(level)

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, err
}
