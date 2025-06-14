package middleware

import (
	"net/http"
	"strings"
)

var allowedOrigins = []string{}

func Cors(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			for _, allowedOrigin := range allowedOrigins {
				if strings.EqualFold(origin, allowedOrigin) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().
						Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
					w.Header().
						Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
					break
				}
			}
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}
