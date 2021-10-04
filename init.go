package main

import "github.com/sirupsen/logrus"
import nested "github.com/antonfisher/nested-logrus-formatter"

func Init() {
	initLog()
}

func initLog() {
	logrus.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		HideKeys:        true,
	})
}
