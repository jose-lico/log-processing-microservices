package logging

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Logger.Info("Created Zap Logger")
}
