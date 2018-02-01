package utils

import (
	"strings"
	"path/filepath"
	"os"
)

//获取路径的父级目录
func DirParent(directory string) string {
	return StringSub(directory, 0, strings.LastIndex(directory, "/"))
}

//获取当前项目路径
func DirCurrent() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//文件或文件夹是否存在
func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}