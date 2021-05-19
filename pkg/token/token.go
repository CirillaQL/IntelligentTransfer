package token

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var MySecret = []byte("IntelligentTransfer")

// MyClaims 自定义的TokenClaims结构体
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
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: ")
		}
		return MySecret, nil
	})
	if token == nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
