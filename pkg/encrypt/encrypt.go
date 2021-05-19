/*
 加密包，用于将用户信息存储在数据库时采用密文格式
*/
package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
)

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

//AesEncrypt 加密函数
func AesEncrypt(input string) (string, error) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	plaintext := []byte(input)
	c := make([]byte, aes.BlockSize+len(plaintext))
	iv := c[:aes.BlockSize]
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	result := base64.StdEncoding.EncodeToString(crypted)
	return result, nil
}

// AesDecrypt 解密函数
func AesDecrypt(input string) (string, error) {
	key, _ := hex.DecodeString("6368616e676520746869732070617373")
	//1.base64解密input
	decodeed, err := base64.StdEncoding.DecodeString(input)
	c := make([]byte, aes.BlockSize+len(decodeed))
	iv := c[:aes.BlockSize]
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(decodeed))
	blockMode.CryptBlocks(origData, decodeed)
	origData = PKCS7UnPadding(origData)
	return string(origData), nil
}
