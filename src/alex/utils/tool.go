// 字符串处理
package utils

import (
	"time"
	"math/rand"
	"github.com/renstrom/shortuuid"
	"github.com/satori/go.uuid"
)


// 生成短字符UUID
func GetShortUUID() string {
	return shortuuid.New()
}

// 生成标准字符UUID
func GetUUID() string {
	return uuid.NewV4().String()
}

// 生成随机字符串
func GetRandString(value int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < value; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// 生成随机数字字符串
func GetRandNum(value int) string {
	str := "123456789"
	bytes := []byte(str)
	result := []byte{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < value; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}