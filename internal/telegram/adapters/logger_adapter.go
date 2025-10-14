package adapters

import "go.uber.org/zap"

type ZapTelegoLogger struct {
	logger *zap.SugaredLogger
}

func NewZapTelegoLogger(logger *zap.SugaredLogger) *ZapTelegoLogger {
	return &ZapTelegoLogger{logger: logger}
}

func (z *ZapTelegoLogger) Debugf(format string, v ...any) {
	z.logger.Debugf(format, v...)
}

func (z *ZapTelegoLogger) Errorf(format string, v ...any) {
	z.logger.Errorf(format, v...)
}
