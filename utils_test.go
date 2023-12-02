package utils

import (
	"testing"
)

func TestFilePath(t *testing.T) {
	filePath, err := NewFilePath()
	if err != nil {
		t.Error("NewFilePath()", err.Error())
		return
	}
	if err = LogServ.InitLog(filePath, &map[string]string{"debuglog": "debug", "errorlog": "error", "infolog": "info"}); err != nil {
		t.Error("InitLog", err.Error())
		return
	}
	if err = LogServ.WriteLog("debuglog", "test for debug"); err != nil {
		t.Error(err.Error())
		return
	}

}
