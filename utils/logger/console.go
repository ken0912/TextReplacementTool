package logger

import (
	"fmt"
	"time"
)

const (
	UNKNOW loglevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

type loglevel uint

type Logger struct {
	Level loglevel
}

//NewConsoleLog ....
func NewConsoleLog(loglevel string) *Logger {
	level, err := transferLogLevel(loglevel)
	if err != nil {
		panic(err)
	}
	return &Logger{
		Level: level,
	}
}
func (l *Logger) isEnable(level loglevel) bool {
	return level >= l.Level
}

func (l *Logger) log(loglevel loglevel, msg string, a ...interface{}) {
	if l.isEnable(loglevel) {
		msg = fmt.Sprintf(msg, a...)
		now := time.Now().Format("2006-01-02 03:04:05:000")
		level := getLogLevelString(loglevel)

		funcName, fileName, line := getInfo(3)
		fmt.Printf("[%v]  [%v]  [%v;%v;%v]  %s\n", now, level, funcName, fileName, line, msg)
	}
}

func (l *Logger) Debug(msg string, a ...interface{}) {
	l.log(DEBUG, msg, a...)
}

func (l *Logger) Trace(msg string, a ...interface{}) {
	l.log(TRACE, msg, a...)
}

func (l *Logger) Info(msg string, a ...interface{}) {
	l.log(INFO, msg, a...)
}

func (l *Logger) Warning(msg string, a ...interface{}) {
	l.log(WARNING, msg, a...)
}

func (l *Logger) Error(msg string, a ...interface{}) {
	l.log(ERROR, msg, a...)
}

func (l *Logger) Fatal(msg string, a ...interface{}) {
	l.log(FATAL, msg, a...)
}
