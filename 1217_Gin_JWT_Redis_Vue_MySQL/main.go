package main

import (
	"os"
	"token/db"
	"token/redis"
	"token/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// MySQL
	err := db.InitMysql()
	if err != nil {
		log.Printf("initMysql() invalid : %s", err.Error())
		return
	}
	defer db.SQLDB.Close()
	// Redis
	err = redis.InitRedis()
	if err != nil {
		log.Printf("InitRedis Error :%s", err.Error())
		return
	}
	defer redis.Client.Close()
	// Router
	err = router.Router()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Router Error")
		return
	}
}
