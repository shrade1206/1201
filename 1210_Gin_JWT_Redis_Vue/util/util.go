package util

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type TokenData struct {
	Token string `json:"token"`
}
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var TokenExpireDuration = time.Minute * 10
var MySecret = []byte("jiyuu")

// 生成token
func GenToken(username string) (string, error) {
	t := MyClaims{
		username, // 自訂Header
		jwt.StandardClaims{ // 設定payload
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "Larry",
		},
	}
	// 選擇編碼模式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	// 用指定的SecretKey加密獲得Token字串
	return token.SignedString(MySecret)
}

// 解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		log.Printf("ParseToken Error :%s", err.Error())
		return nil, err
	}
	// 驗證claims正確就回傳
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid Token")
}
