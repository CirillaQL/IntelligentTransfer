package token

import (
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("IntelligentTransfer")

//自定义的TokenClaims结构体
type MyClaims struct {
	UUid        string `json:"UUid"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

// GenToken 生成对应的Token
func GenToken(UUid, PhoneNumber string) (string, error) {
	c := MyClaims{
		UUid:        UUid,
		PhoneNumber: PhoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(MySecret)
}

// ParseToken 解析Token
