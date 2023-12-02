package utils

import (
	"fmt"
	"os"
	"strings"
)

type TFilepath struct {
	CurrentPath string
	DirFlag     string
	FileDirs    *map[string]string
}

// NewFilePath get current path and os dir flag
func NewFilePath() (*TFilepath, error) {
	currentPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dirFlag := "/"
	if strings.Contains(currentPath, "\\") {
		dirFlag = "\\"
	}
	var enStr TEnString
	enStr.Load(currentPath)
	currentPath = enStr.CutFromLast(dirFlag) + dirFlag

	return &TFilepath{currentPath, dirFlag, nil}, nil
}

func (fp *TFilepath) SetFileDir(fileDirs *map[string]string) error {
	var err error
	var fullPath string
	checkFilePath := func(filePath string) error {
		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			err = os.Mkdir(filePath, 0766)
			if err != nil {
				return fmt.Errorf("创建目录%s出错:%s", filePath, err.Error())
			}
		}
		return nil
	}
	for key, val := range *fileDirs {
		fullPath = fp.CurrentPath + val + fp.DirFlag
		if err = checkFilePath(fullPath); err != nil {
			return err
		}
		(*fileDirs)[key] = fullPath
	}
	fp.FileDirs = fileDirs
	return nil

}
func (fp *TFilepath) GetFileDir(fileType string) (string, error) {
	var fullPath string
	var ok bool
	if fullPath, ok = (*fp.FileDirs)[fileType]; ok {
		return fullPath, nil
	}
	return "", fmt.Errorf("%s not exists", fileType)
}
