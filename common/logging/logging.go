package logging

import "go.uber.org/zap"

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer Logger.Sync()
}
