package main

import (
	"github.com/UndeadDemidov/ya-pr-diploma/internal/server"
)

func main() {
	// Warning! init() is in use!
	srv, _ := server.NewServer(cfg.Server)
	srv.Run()
}