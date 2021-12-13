package session

import (
	"log"

	"github.com/gin-contrib/sessions/redis"
)

var Store redis.Store

func Initsession() (err error) {
	// 使用 "github.com/gin-contrib/sessions/redis" 跟redis連接
	Store, err = redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Printf("Session Error :%s", err.Error())
		return err
	}
	return
}
