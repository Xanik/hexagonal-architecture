package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var DefaultLogger = logrus.New()

func init() {
	DefaultLogger.Out = io.MultiWriter(os.Stdout)
	DefaultLogger.SetLevel(logrus.TraceLevel)
	DefaultLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
}
