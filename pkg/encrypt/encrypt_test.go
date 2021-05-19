package encrypt

import (
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {
	plaintext := "这是一个测试"
	fmt.Println("原文本: " + plaintext)
	ans, err := AesEncrypt(plaintext)
	fmt.Println("加密结果: " + ans)
	if err != nil {
		t.Error(err)
	}
	fmt.Print("解密结果: ")
	result, err := AesDecrypt(ans)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(result)
}
