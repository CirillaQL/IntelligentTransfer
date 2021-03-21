/*
 加密包，用于将用户信息存储在数据库时采用密文格式
*/
package encrypt

import (
	"IntelligentTransfer/pkg/logger"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

//加密
func Encrypt(input string) (string, error) {
	if input == "" {
		logger.Error("encrypt input is empty")
		return "", fmt.Errorf("encrypt input is empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("generate failed err: %+v", err)
		return "", err
	}
	return string(hash), nil
}

//判断是否相等
func CheckEncryptString(origin, encryption string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(encryption), []byte(origin))
	if err != nil {
		return false
	}
	return true
}
