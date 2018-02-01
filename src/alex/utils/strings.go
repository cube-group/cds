package utils

import (
	"fmt"
	"strings"
)

//截取字符串
func StringSub(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

//不详内容合并为字符串
func StringJoin(sep string, values ...interface{}) string {
	var arr []string
	for _, item := range values {
		arr = append(arr, fmt.Sprint(item))
	}
	return strings.Join(arr, sep)
}

//不详内容合并为打印字符串
func StringTrace(values ...interface{}) string {
	return fmt.Sprint(values...)
}