package main

import (
	"TT-Watch/cmd"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.99999",
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	cmd.Watch.Execute()
}
