package logger

import (
	"errors"
	"path"
	"runtime"
	"strings"
)

func transferLogLevel(level string) (loglevel, error) {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return DEBUG, nil
	case "TRACE":
		return TRACE, nil
	case "INFO":
		return INFO, nil
	case "WARNING":
		return WARNING, nil
	case "ERROR":
		return ERROR, nil
	case "FATAL":
		return FATAL, nil
	default:
		err := errors.New("invalid log level!")
		return UNKNOW, err
	}
}

func getInfo(n int) (funcName, filename string, line int) {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		panic("getInfo error:")
	}
	funcName = runtime.FuncForPC(pc).Name()
	filename = path.Base(file)
	return funcName, file, line
}

func getLogLevelString(loglevel loglevel) string {
	switch loglevel {
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOW"
	}
}
