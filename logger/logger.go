package logger

import (
	"io"
	"os"
	"qwerty/config"

	log "github.com/sirupsen/logrus"
)

const (
	toStdOut = iota
	toFile
	both
)

func getLogLevel(level string) log.Level {
	switch level {
	case "DEBUG":
		return log.DebugLevel
	case "ERROR":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

func Setup(conf *config.Logger) {
	//defer f.Close()

	if conf.Format == "json" {
		log.SetFormatter(&log.JSONFormatter{TimestampFormat: conf.TimeStampFormat})
	} else {
		format := &log.TextFormatter{DisableColors: true, TimestampFormat: conf.TimeStampFormat}
		log.SetFormatter(format)
	}
	if conf.Log2Engine == toStdOut {
		log.SetOutput(os.Stdout)
	} else {
		f, err := os.OpenFile(conf.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("Error opening log file: ", err)
		}
		if conf.Log2Engine == toFile {
			log.SetOutput(f)
		} else if conf.Log2Engine == both {
			mw := io.MultiWriter(os.Stdout, f)
			log.SetOutput(mw)
		}
	}
	loglevel := getLogLevel(conf.Level)
	log.SetLevel(loglevel)
}
