package utils

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	ContentTypeKey  = "Content-Type"
	ContentTypeJSON = "application/json"
	ContentTypeText = "text/plain"
)

func InternalServerError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Error().Err(err).Msg("")
}

func ServerError(w http.ResponseWriter, err error, status int) {
	http.Error(w, err.Error(), status)
	log.Error().Err(err).Msg("")
}

func TimeParseHelper(layout string, t string) time.Time {
	tmp, err := time.Parse(layout, t)
	if err != nil {
		panic(err)
	}
	return tmp
}

func TimeRFC3339ParseHelper(t string) time.Time {
	return TimeParseHelper(time.RFC3339, t)
}
