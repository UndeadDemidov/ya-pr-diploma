package middleware

import (
	"compress/gzip"
	"net/http"

	"github.com/UndeadDemidov/ya-pr-diploma/internal/presenter/http/utils"
)

// Decompress реализует распаковку запроса переданного в сжатом gzip
func Decompress(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				utils.InternalServerError(w, err)
			}
			defer func() {
				err := gz.Close()
				if err != nil {
					utils.InternalServerError(w, err)
				}
			}()
			r.Body = gz
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
