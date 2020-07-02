package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

type LogMsg struct {
	TimeStamp string
	Level     loglevel
	FuncName  string
	FileName  string
	LineNo    int
	Msg       string
}

type FileLogger struct {
	Level        loglevel
	FilePath     string
	FileName     string
	FileObj      *os.File
	ErrorFileObj *os.File
	FileMaxSize  int
	LogChan      chan *LogMsg
}

//NewFileLog ....
func NewFileLog(loglevel, fp, fn string, fms int) *FileLogger {
	level, err := transferLogLevel(loglevel)
	if err != nil {
		panic(err)
	}
	fl := &FileLogger{
		Level:       level,
		FilePath:    fp,
		FileName:    fn,
		FileMaxSize: fms,
		LogChan:     make(chan *LogMsg, 50000),
	}
	fl.initFile()
	return fl
}

//init filelog
func (fl *FileLogger) initFile() error {
	fullname := path.Join(fl.FilePath, fl.FileName)
	fileobj, err := os.OpenFile(fullname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("initFile error:", err)
		return err
	}
	fl.FileObj = fileobj

	errFileObj, err := os.OpenFile("Err"+fullname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("initFile error:", err)
		return err
	}
	fl.ErrorFileObj = errFileObj
	// go fl.writeLogToFile()
	return nil
}

func (fl *FileLogger) isEnable(level loglevel) bool {
	return level >= fl.Level
}

func (fl *FileLogger) checkFileSize(fileobj *os.File) bool {
	fileInfo, err := fileobj.Stat()
	if err != nil {
		fmt.Println("CheckFileSize error:", err)
	}
	fileSize := fileInfo.Size()
	if fileSize >= int64(fl.FileMaxSize) {
		return true
	}
	return false
}

func (fl *FileLogger) splitFileLog(fileobj *os.File) (*os.File, error) {
	nowstr := time.Now().Format("20060102030405000")
	fileinfo, err := fileobj.Stat()
	if err != nil {
		fmt.Println("get file info error:", err)
		return nil, err
	}
	logname := path.Join(fl.FilePath, fileinfo.Name())
	//1 close current file
	fileobj.Close()
	//2 backup current file
	os.Rename(logname, logname+nowstr)
	//3 open new file
	fileobj, err = os.OpenFile(logname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("log() error:", err)
		return nil, err
	}
	//4 assign the new fileobj to original fl.FileObj
	return fileobj, nil
}

/*
func (fl *FileLogger) writeLogToFile() {
	for {
		if fl.checkFileSize(fl.FileObj) {
			fileobj, err := fl.splitFileLog(fl.FileObj)
			if err != nil {
				return
			}
			fl.FileObj = fileobj
		}
		select {
		case logmsg := <-fl.LogChan:
			msg := fmt.Sprintf("[%v]  [%v]  [%v;%v;%v]  %s\r\n", logmsg.TimeStamp, getLogLevelString(logmsg.Level), logmsg.FuncName, logmsg.FileName, logmsg.LineNo, logmsg.Msg)
			fmt.Fprintf(fl.FileObj, msg)

			if logmsg.Level >= ERROR {
				if fl.checkFileSize(fl.ErrorFileObj) {
					fileobj, err := fl.splitFileLog(fl.ErrorFileObj)
					if err != nil {
						return
					}
					fl.ErrorFileObj = fileobj
				}
				fmt.Fprintf(fl.ErrorFileObj, msg)
			}
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}

}
*/
func (fl *FileLogger) log(loglevel loglevel, msg string, a ...interface{}) {
	if fl.isEnable(loglevel) {
		msg = fmt.Sprintf(msg, a...)
		now := time.Now().Format("2006-01-02 03:04:05:000")

		funcName, fileName, line := getInfo(3)
		logMsg := &LogMsg{
			TimeStamp: now,
			Level:     loglevel,
			FuncName:  funcName,
			FileName:  fileName,
			LineNo:    line,
			Msg:       msg,
		}
		// fl.LogChan <- logMsg

		msg := fmt.Sprintf("[%v]  [%v]  [%v;%v;%v]  %s\r\n", logMsg.TimeStamp, getLogLevelString(logMsg.Level), logMsg.FuncName, logMsg.FileName, logMsg.LineNo, logMsg.Msg)
		if logMsg.Level == 2 || logMsg.Level == 3 {
			msg = fmt.Sprintf("[%v]  [%v]  %s\r\n", logMsg.TimeStamp, getLogLevelString(logMsg.Level), logMsg.Msg)
		}
		fmt.Fprintf(fl.FileObj, msg)

		if logMsg.Level >= ERROR {
			if fl.checkFileSize(fl.ErrorFileObj) {
				fileobj, err := fl.splitFileLog(fl.ErrorFileObj)
				if err != nil {
					return
				}
				fl.ErrorFileObj = fileobj
			}
			fmt.Fprintf(fl.ErrorFileObj, msg)
		}

	}
}

func (fl *FileLogger) Debug(msg string, a ...interface{}) {
	fl.log(DEBUG, msg, a...)
}

func (fl *FileLogger) Trace(msg string, a ...interface{}) {
	fl.log(TRACE, msg, a...)
}

func (fl *FileLogger) Info(msg string, a ...interface{}) {
	fl.log(INFO, msg, a...)
}

func (fl *FileLogger) Warning(msg string, a ...interface{}) {
	fl.log(WARNING, msg, a...)
}

func (fl *FileLogger) Error(msg string, a ...interface{}) {
	fl.log(ERROR, msg, a...)
}

func (fl *FileLogger) Fatal(msg string, a ...interface{}) {
	fl.log(FATAL, msg, a...)
}
