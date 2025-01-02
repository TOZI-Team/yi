package devlog

import (
	"github.com/sirupsen/logrus"
)

var DevLog = logrus.New()

func init() {
	DevLog.SetLevel(logrus.DebugLevel)
}
