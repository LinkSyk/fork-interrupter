package xlog

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = New()

type Xlog struct {
	l *logrus.Logger
}

func setLogLevel(l *logrus.Logger) {
	v := os.Getenv("LOGLEVEL")
	switch v {
	case "info":
		l.SetLevel(logrus.InfoLevel)
	case "debug":
		l.SetLevel(logrus.DebugLevel)
	case "warn":
		l.SetLevel(logrus.WarnLevel)
	}
}

func New() *Xlog {
	l := logrus.New()
	setLogLevel(l)
	return &Xlog{
		l: l,
	}
}

func Infof(fmt string, args ...any) {
	log.l.Infof(fmt, args...)
}

func Infoln(args ...any) {
	log.l.Infoln(args...)
}

func Debugf(fmt string, args ...any) {
	log.l.Debugf(fmt, args...)
}

func Debugln(args ...any) {
	log.l.Debugln(args...)
}

func Fatalf(fmt string, args ...any) {
	log.l.Fatalf(fmt, args...)
}

func Fatalln(args ...any) {
	log.l.Fatalln(args...)
}
