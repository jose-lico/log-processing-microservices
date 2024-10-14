package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger(env string) {
	var err error

	if env == "LOCAL" {
		Logger, err = zap.NewDevelopment()
	} else {
		Logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}
	Logger.Info("Created Zap Logger")
}
