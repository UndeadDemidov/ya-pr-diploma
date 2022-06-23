package utils

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

const (
	ContentTypeKey  = "Content-Type"
	ContentTypeJSON = "application/json"
)

func InternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Error().Err(err)
}

func ServerError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
	log.Error().Err(err)
}
