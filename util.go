package util

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type TokenData struct {
	Token string `json:"token"`
}

const TokenExpire = time.Hour * 2

var MySecret = []byte("jiyuu")

// 生成Token
func GenToken(username string) (string, error) {
	t := MyClaims{
		username, //自訂Header
		jwt.StandardClaims{ //設定Playload
			ExpiresAt: time.Now().Add(TokenExpire).Unix(),
			Issuer:    "Larry",
		},
	}
	// 選擇編碼模式
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	// 之後需要補上Token存進資料庫，確保登錄的時候Token還在
	tk, err := token.SignedString(MySecret)
	if err != nil {
		log.Printf("token Error: %s", err.Error())
		return "nil", err
	}
	logtoken := TokenData{Token: tk}
	fmt.Println(logtoken)
	// 用指定的SecretKey來穫得編碼後的字串
	return token.SignedString(MySecret)
}

//解析Token
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		log.Printf("ParseToken Error :%s", err.Error())
		return nil, err
	}
	// 驗證token是否正確
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
