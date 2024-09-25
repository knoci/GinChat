package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写，计算字符串的 MD5 哈希值
func Md5Encode(data string) string {
	// 创建了一个新的 MD5 哈希计算器
	h := md5.New()
	// 将输入的字符串 data 转换为字节切片（[]byte(data)），然后写入到 MD5 哈希计算器 h 中
	h.Write([]byte(data))
	// h.Sum 方法调用会完成哈希计算，并将结果返回为一个字节切片。nil 参数表示没有额外的字节切片需要追加到哈希计算中
	tempStr := h.Sum(nil)
	// 将上一步得到的字节切片 tempStr 转换为十六进制字符串
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 随机数加密
func MakePassword(plainpwd, salt string) string {
	return Md5Encode(plainpwd + salt)
}

// 随机数解密
func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
