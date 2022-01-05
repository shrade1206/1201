package main

import (
	"todoList/db"
	"todoList/redis"
	"todoList/router"

	"github.com/rs/zerolog/log"
)

func main() {
	// MySQL
	err := db.InitMysql()
	if err != nil {
		log.Fatal().Err(err).Msg("InitMysql invalid")
		return
	}
	defer db.SQLDB.Close()
	// Redis
	err = redis.InitRedis()
	if err != nil {
		log.Fatal().Err(err).Msg("InitRedis invalid")
		return
	}
	defer redis.Client.Close()
	// Router
	err = router.Router()
	if err != nil {
		log.Fatal().Err(err).Msg("Router invalid")
		return
	}

}
