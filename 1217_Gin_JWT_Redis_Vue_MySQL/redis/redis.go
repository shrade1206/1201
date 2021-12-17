package redis

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis"
)

type Dsn struct {
	Addr        string
	Password    string
	DB          int
	PoolSize    int
	PoolTimeout int
	MaxConnAge  int
}

var Client *redis.Client

func InitRedis() (err error) {
	// 讀取Redis_Config，並寫入
	var dsn Dsn
	file, err := os.Open("./config/Redis_Config.json")
	if err != nil {
		return
	}
	data := json.NewDecoder(file)
	err = data.Decode(&dsn)
	if err != nil {
		return
	}
	// dsn, err := util.RedisDsn("./config/Redis_Config.json")
	//連線Redis
	Client = redis.NewClient(&redis.Options{
		Addr:        dsn.Addr,
		Password:    dsn.Password,
		DB:          dsn.DB,
		PoolSize:    dsn.PoolSize,
		PoolTimeout: time.Duration(dsn.PoolTimeout) * time.Second,
		MaxConnAge:  time.Duration(dsn.MaxConnAge) * time.Second,
	})
	//確認是否連到Redis
	err = Client.Ping().Err()
	if err != nil {
		log.Printf("Redis Ping Error :%s", err.Error())
		return err
	}
	return
}
