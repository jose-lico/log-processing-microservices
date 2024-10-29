package kafka

import (
	"fmt"

	"go.uber.org/zap"
)

type ZapSaramaLogger struct {
	logger *zap.Logger
}

func (l *ZapSaramaLogger) Print(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}

func (l *ZapSaramaLogger) Printf(format string, v ...interface{}) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *ZapSaramaLogger) Println(v ...interface{}) {
	l.logger.Info(fmt.Sprintln(v...))
}
