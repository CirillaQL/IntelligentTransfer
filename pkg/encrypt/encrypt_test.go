package encrypt

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	ans, err := Encrypt("sdjisoad")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ans)
}

func TestCheckEncryptString(t *testing.T) {
	passwordOld := "test if pass"
	passwordNex, err := Encrypt(passwordOld)
	if err != nil {
		t.Error(err)
	}
	ans := CheckEncryptString(passwordOld, passwordNex)
	if ans == false {
		t.Error("wrong answer")
	}
}
