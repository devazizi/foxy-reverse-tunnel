package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func Debug(section string, msg string, fields map[string]interface{}) {
	loggerPkg := logrus.New()
	loggerPkg.Out = os.Stdout
	loggerPkg.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	loggerPkg.WithFields(fields).Info(msg)
}

func Error(section string, msg string, fields map[string]interface{}) {
	loggerPkg := logrus.New()
	loggerPkg.Out = os.Stdout
	loggerPkg.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	loggerPkg.WithFields(fields).Error(msg)
}
