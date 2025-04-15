package router

import (
	"chirpy/config"
	"chirpy/handlers"
	"chirpy/middleware"
	"net/http"
)

func New(apiConfig *config.Configuration) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/app/", middleware.MetricsInc(apiConfig, http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlers.HandleHealthz)

	mux.HandleFunc("GET /admin/metrics", handlers.HandleAdminMetrics(apiConfig))
	mux.HandleFunc("POST /admin/reset", handlers.HandleAdminReset(apiConfig))

	mux.HandleFunc("POST /api/users", handlers.HandleCreate(apiConfig))
	mux.HandleFunc("POST /api/login", handlers.HandleLogin(apiConfig))
	mux.HandleFunc("POST /api/refresh", handlers.HandleRefresh(apiConfig))
	mux.HandleFunc("POST /api/revoke", handlers.HandleRevoke(apiConfig))

	mux.HandleFunc("POST /api/chirps", handlers.HandleCreateChirp(apiConfig))
	mux.HandleFunc("GET /api/chirps", handlers.HandleListChirps(apiConfig))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.HandleGetChirp(apiConfig))

	return mux
}
