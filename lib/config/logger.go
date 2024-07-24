package config

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {
	logger := &Logger{
		log: logrus.New(),
	}
	logger.init()
	return logger
}

func (l *Logger) init() {
	// Open a file for writing logs
	file, err := os.OpenFile("mazimart.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Set the output of the logger to the file
		l.log.Out = file
	} else {
		l.log.Info("Failed to log to file, using default stderr")
	}

	l.log.SetFormatter(&logrus.JSONFormatter{})
}

func (l *Logger) Log() *logrus.Logger {
	return l.log
}
