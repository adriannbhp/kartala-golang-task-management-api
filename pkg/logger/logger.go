package logger

import "github.com/sirupsen/logrus"

var Logger = logrus.New()

func SetupLogger() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.Info("Logger initialized!")
	Logger = log
}
