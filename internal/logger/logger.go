package logger

import (
	"io"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// New creates a logrus logger writing to stdout and rotating file.
func New(logFile string, level logrus.Level) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(level)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	if logFile != "" {
		log.SetOutput(io.MultiWriter(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    100, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   false,
		}, os.Stdout))
	} else {
		log.SetOutput(os.Stdout)
	}

	return log
}
