package middleware

import (
	"chirpy/config"
	"net/http"
)

func MetricsInc(cfg *config.Configuration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
