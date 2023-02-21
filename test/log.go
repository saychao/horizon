package test

import (
	"github.com/saychao/horizon/log"
	"github.com/sirupsen/logrus"
)

var testLogger *log.Entry

func init() {
	testLogger, _ = log.New()
	testLogger.Entry.Formatter(&logrus.TextFormatter{DisableColors: true})
	testLogger.Entry.Level(log.DebugLevel)
}
