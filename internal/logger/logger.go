package logger

import (
	"github.com/sirupsen/logrus"
	// "os"
)

type Logger struct {
	*logrus.Logger
}

func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	// 	logger.Fatalf("Error opening file: %v", err)
	// }

	// logger.SetOutput(file)

	logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			DisableColors: false,
			FullTimestamp: true,
	})

	logger.SetLevel(logrus.InfoLevel)

	return logger
}
