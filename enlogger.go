package utils

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func init() {
	logger := make(map[string]*EnLogger)
	LogServ = TLogService{logger}
}

var LogServ TLogService

type EnLogger struct {
	logger   *log.Logger
	logFile  *os.File
	filePath string
	fileName string
	lock     *sync.Mutex
}
type TLogService struct {
	Logger map[string]*EnLogger
}

func NewLogger(filePath string) (*EnLogger, error) {
	fileName := filePath + "log_" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return nil, err
	}
	loger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	loger.SetOutput(logFile)
	lock := new(sync.Mutex)
	return &EnLogger{loger, logFile, filePath, fileName, lock}, nil
}
func (enLog *EnLogger) CloseLog() {
	_ = enLog.logFile.Close()
}
func (enLog *EnLogger) newFile() {
	_ = enLog.logFile.Close()
	enLog.fileName = enLog.filePath + "log_" + time.Now().Format("20060102") + ".log"
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

func (ls *TLogService) InitLog(fp *TFilepath, logs *map[string]string) error {
	var err error
	if err = fp.SetFileDir(logs); err != nil {
		return err
	}
	for key, val := range *fp.FileDirs {
		if ls.Logger[key], err = NewLogger(val); err != nil {
			return err
		}
	}
	return nil
}
func (ls *TLogService) WriteLog(logType string, v ...any) error {
	var logWriter *EnLogger
	var ok bool
	if logWriter, ok = ls.Logger[logType]; !ok {
		return fmt.Errorf("%s not exists", logType)
	}
	logWriter.WriteLog(v...)
	return nil
}
