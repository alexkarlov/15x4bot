package log

import (
	"log"
)

type LogLevel int

const (
	LOG_LEVEL_ERROR LogLevel = 0
	LOG_LEVEL_WARN  LogLevel = 1
	LOG_LEVEL_DEBUG LogLevel = 2
)

var level = LOG_LEVEL_ERROR

// SetLevel sets the logger level. Default level value = LOG_LEVEL_ERROR
func SetLevel(l LogLevel) {
	level = l
}

func Infof(f string, msg ...interface{}) {
	if level != LOG_LEVEL_DEBUG {
		return
	}
	log.Printf(f, msg...)
}

func Info(msg ...interface{}) {
	if level != LOG_LEVEL_DEBUG {
		return
	}
	log.Print(msg...)
}

func Warnf(f string, msg ...interface{}) {
	if level == LOG_LEVEL_ERROR {
		return
	}
	log.Printf(f, msg...)
}

func Errorf(f string, msg ...interface{}) {
	log.Printf(f, msg...)
}

func Error(msg ...interface{}) {
	log.Print(msg...)
}
