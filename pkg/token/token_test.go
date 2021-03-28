package token

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	ans, err := GenToken("sdadasd", "15840613358")
	if err != nil {
		t.Error(err)
	}
	a, _ := ParseToken(ans)
	fmt.Println(a.UUid)
}
