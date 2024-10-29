package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger(env string) {
	var err error
	var cfg zap.Config

	if env == "LOCAL" {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.StacktraceKey = ""
		Logger, err = cfg.Build()
	} else {
		cfg = zap.NewProductionConfig()
		Logger, err = cfg.Build()
	}

	if err != nil {
		panic(err)
	}
	Logger.Info("Created Zap Logger")
}
