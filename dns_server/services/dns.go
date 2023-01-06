package services

import (
	//"fmt"

	"github.com/sirupsen/logrus"
)

func TT(l *logrus.Entry) {
	// tlog := l.WithField("func", "TT")
	// tlog.Info("++++++ tt ++++++++")
	l.Info("++++++ tt ++++++++")
}
