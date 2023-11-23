package utils

import (
	"fmt"
	"os"
)

type TFilepath struct {
	CurrentPath string
	DirFlag     string
	FilePaths   map[string]string
}

// InitFilePath 初始化文件目录,当前目录开始,目录分隔符为"/"
func (fp *TFilepath) InitFilePath(parentPath, dir string, filePaths *map[string]string) error {
	var err error
	fp.CurrentPath = parentPath
	fp.DirFlag = dir
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
	fp.FilePaths = make(map[string]string)

	for key, val := range *filePaths {
		if err = checkFilePath(val); err != nil {
			return err
		}
		fp.FilePaths[key] = val
	}
	return nil

}
