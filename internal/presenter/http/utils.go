package http

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

func InternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Error().Err(err)
}
