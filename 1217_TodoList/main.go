package main

import (
	"os"
	"todoList/db"
	"todoList/router"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// MySQL
	err := db.InitMysql()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("MySQL Error")
		return
	}
	defer db.SQLDB.Close()
	// Router
	err = router.Router()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Router Error")

		return
	}

}
