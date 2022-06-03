package api

import (
	"net/http"
	"strings"
)

func Dispatcher(
	ping *PingHandler,
	promotion *PromotionHandler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if strings.HasPrefix(r.RequestURI, "/promotions") {
			promotion.GetByID(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/ping") {
			ping.Ping(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	})
}
