package main

import (
	"github.com/KihaRaito/sofupo-backend/backend"
	"github.com/rs/zerolog/log"
	// log "github.com/sirupsen/logrus"
)

func main() {
	log.Logger = log.With().Caller().Logger()

	// databaseの初期化
	log.Info().Msg("initializing database ...")
	db := backend.DbConn()
	backend.DbInit(db)

	log.Info().Msg("server listening ...")
	backend.Router().Run(":8080")
}
