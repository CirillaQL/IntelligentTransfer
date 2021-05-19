package token

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	fmt.Println("Token信息: UUid: 555c40c8-f8b2-49a3-b82d-12b861ea8429  电话号码: 158xxxxxxxx")
	ans, err := GenToken("555c40c8-f8b2-49a3-b82d-12b861ea8429", "158xxxxxxxx")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Token: " + ans)
	a, _ := ParseToken(ans)
	fmt.Print("解析Token信息: UUid: ")
	fmt.Print(a.UUid)
	fmt.Print("   电话号码: ")
	fmt.Println(a.PhoneNumber)
}
