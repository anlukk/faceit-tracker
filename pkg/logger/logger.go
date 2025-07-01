package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

var (
	logger      *zap.Logger
	sugar       *zap.SugaredLogger
	once        sync.Once
	atomicLevel zap.AtomicLevel
)

func BuildLogger(logLevel string) (*zap.Logger, *zap.SugaredLogger, error) {
	var initErr error
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevel()
		if err := setLogLevel(logLevel); err != nil {
			initErr = err
			return
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atomicLevel,
		)

		logger = zap.New(
			core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
		sugar = logger.Sugar()
	})

	if initErr != nil {
		return nil, nil, initErr
	}
	return logger, sugar, nil
}

func GetSugaredLogger() *zap.SugaredLogger {
	if sugar == nil {
		panic("logger not initialized, call BuildLogger first")
	}
	return sugar
}

func setLogLevel(logLevel string) error {
	switch strings.ToUpper(logLevel) {
	case LevelDebug:
		atomicLevel.SetLevel(zapcore.DebugLevel)
	case LevelInfo:
		atomicLevel.SetLevel(zapcore.InfoLevel)
	case LevelWarn:
		atomicLevel.SetLevel(zapcore.WarnLevel)
	case LevelError:
		atomicLevel.SetLevel(zapcore.ErrorLevel)
	default:
		return fmt.Errorf("invalid log level: must be DEBUG, INFO, WARN or ERROR")
	}
	return nil
}

func CurrentLevel() string {
	return atomicLevel.String()
}
