package main

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/server"
	"github.com/rs/zerolog/log"
)

func main() {
	// Warning! init() is in use!
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("server is not configured")
	}
	srv.Run()
}
