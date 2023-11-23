package utils

import "fmt"

type TLogService struct {
	CurrentPath string
	DirFlag     string
	ErrLog      *EnLogger
	InfoLog     *EnLogger
	DebugLog    *EnLogger
	//SysCfg      *TServerCfg
}

var LogServ TLogService

func (env *TLogService) InitLog(infoLog, errlog, debuglog string) error {
	var err error

	infoLog = fmt.Sprintf("%s%s%s", env.CurrentPath, infoLog, env.DirFlag)
	errlog = fmt.Sprintf("%s%s%s", env.CurrentPath, errlog, env.DirFlag)
	debuglog = fmt.Sprintf("%s%s%s", env.CurrentPath, debuglog, env.DirFlag)
	var fp TFilepath
	files := map[string]string{"infolog": infoLog, "errorlog": errlog, "debuglog": debuglog}
	err = fp.InitFilePath(env.CurrentPath, env.DirFlag, &files)
	if err != nil {
		return err
	}

	if env.InfoLog, err = NewLogger(infoLog, InfoLog); err != nil {
		return err
	}
	if env.ErrLog, err = NewLogger(errlog, ErrorLog); err != nil {
		return err
	}
	if env.DebugLog, err = NewLogger(debuglog, DebugLog); err != nil {
		return err
	}
	return nil
}

func (env *TLogService) GetFilePath(fileName, filePath string) string {
	fullPath := filePath + fileName
	if fullPath == "" {
		fullPath = env.CurrentPath + fileName
	}
	var enStr TEnString
	enStr.String = fullPath
	//确保右侧有分隔符
	return enStr.TrimFromRight(env.DirFlag) + env.DirFlag
}
