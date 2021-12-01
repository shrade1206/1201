package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type MyClaims struct {
	Username string `json:username`
	jwt.StandardClaims
}

func main() {
	mySigningKey := []byte("hiroyukisawano")
	// StandardClaims
	// c := MyClaims{
	// 	Username: "jiyuu",
	// 	StandardClaims: jwt.StandardClaims{
	// 		NotBefore: time.Now().Unix() - 60,      // 什麼時候開始，當前時間-60秒
	// 		ExpiresAt: time.Now().Unix() + 60*60*2, // 什麼時候過期，當前時間+2個小時
	// 		Issuer:    "jiyuu",                     // 誰簽發的
	// 	},
	// }

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":      time.Now().Unix() - 5, // 什麼時候開始
		"exp":      time.Now().Unix() + 5, // 什麼時候到期
		"iss":      time.Now().Unix(),     // 誰簽發的
		"username": "J",
	}) // 使用NewWithClaims，選擇需要的加密方式，加上資料，生成token
	d, err := t.SignedString(mySigningKey) // 使用SignedString把資料加密
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println(d)

	token, err := jwt.ParseWithClaims(d, &MyClaims{}, func(t *jwt.Token) (interface{}, error) {
		return mySigningKey, nil // 回傳要解析的Key
	})
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Println(token.Claims.(*MyClaims).Username)
}
