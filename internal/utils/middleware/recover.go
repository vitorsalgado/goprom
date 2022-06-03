package middleware

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// Recovery is a recover() middleware for all unhandled errors
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovery := recover(); recovery != nil {
				var msg string
				switch x := recovery.(type) {
				case string:
					msg = x
				case error:
					msg = x.Error()
				default:
					msg = http.StatusText(http.StatusInternalServerError)
				}

				log.Error().Timestamp().Stack().
					Interface("recover", recovery).
					Msg(msg)

				w.Header().Set("content-type", "text/plain")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf("an unexpected error occured. %v", msg)))
			}

		}()

		next.ServeHTTP(w, r)
	})
}
