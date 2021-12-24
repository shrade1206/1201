package main

import (
	"token/db"
	"token/redis"
	"token/router"

	"github.com/rs/zerolog/log"
)

func main() {
	// MySQL
	err := db.InitMysql()
	if err != nil {
		log.Fatal().Err(err).Msg("InitRouter invalid")
		return
	}
	defer db.SQLDB.Close()
	// Redis
	err = redis.InitRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("InitRouter invalid")
		return
	}
	defer redis.Client.Close()
	// Router
	err = router.Router()
	if err != nil {
		log.Fatal().Err(err).Msg("InitRouter invalid")
		return
	}
}
