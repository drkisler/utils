package utils

import (
	"log"
	"os"
	"sync"
	"time"
)

type LogType uint8

const (
	InfoLog LogType = iota
	ErrorLog
	DebugLog
)

type ILogger interface {
	WriteLog(v ...interface{})
}
type EnLogger struct {
	logger   *log.Logger
	logFile  *os.File
	filePath string
	fileName string
	lock     *sync.Mutex
	postfix  string
}

func NewLogger(filePath string, logtype LogType) (*EnLogger, error) {
	var fix string
	switch logtype {
	case InfoLog:
		fix = "_info.log"
	case ErrorLog:
		fix = "_error.log"
	case DebugLog:
		fix = "_debug.log"
	}
	fileName := filePath + time.Now().Format("20060102") + fix
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return nil, err
	}
	loger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	loger.SetOutput(logFile)
	lock := new(sync.Mutex)
	return &EnLogger{loger, logFile, filePath, fileName, lock, fix}, nil
}
func (enLog *EnLogger) CloseLog() {
	_ = enLog.logFile.Close()
}
func (enLog *EnLogger) newFile() {
	_ = enLog.logFile.Close()
	enLog.fileName = enLog.filePath + time.Now().Format("20060102") + enLog.postfix
	enLog.logFile, _ = os.OpenFile(enLog.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	enLog.logger.SetOutput(enLog.logFile)
}

func (enLog *EnLogger) WriteLog(v ...interface{}) {
	enLog.lock.Lock()
	defer enLog.lock.Unlock()
	if enLog.fileName != time.Now().Format("20060102") {
		enLog.newFile()
	}
	enLog.logger.Println(v...)

}
