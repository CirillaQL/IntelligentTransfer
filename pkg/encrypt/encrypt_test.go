package encrypt

import (
	"fmt"
	"testing"
)

func TestAesEncrypt(t *testing.T) {

	plaintext := "hello ming"

	ans, err := AesEncrypt(plaintext)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(AesDecrypt(ans))
}
