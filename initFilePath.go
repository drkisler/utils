package utils

import (
	"fmt"
	"os"
	"strings"
)

type TFilepath struct {
	CurrentPath string
	DirFlag     string
	FilePaths   map[string]string
}

// InitFilePath 初始化文件目录,当前目录开始,目录分隔符为"/"
func (fp *TFilepath) InitFilePath(parentPath, dirFlag string, filePaths *map[string]string) error {
	var err error
	arr := strings.Split(parentPath, dirFlag)
	if arr[len(arr)-1] == "" {
		fp.CurrentPath = parentPath
	} else {
		fp.CurrentPath = parentPath + dirFlag
	}
	fp.DirFlag = dirFlag
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
		if err = checkFilePath(fp.CurrentPath + val); err != nil {
			return err
		}
		fp.FilePaths[key] = fp.CurrentPath + val
	}
	return nil

}
