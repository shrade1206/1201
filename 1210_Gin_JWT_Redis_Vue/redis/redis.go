package redis

import (
	"log"
	"time"

	"github.com/go-redis/redis"
)

var Client *redis.Client

func InitRedis() (err error) {

	Client = redis.NewClient(&redis.Options{
		Addr:        "localhost:6379",
		Password:    "", //不設定密碼
		DB:          0,  //預設資料庫
		PoolSize:    10,
		PoolTimeout: time.Hour * 1,
		MaxConnAge:  time.Hour * 1,
	})
	//確認是否連到Redis
	_, err = Client.Ping().Result()
	if err != nil {
		log.Printf("Redis Ping Error :%s", err.Error())
		return err
	}
	return
}
