package main

import (
	"chirpy/config"
	"chirpy/handlers"
	"chirpy/internal/database"
	"chirpy/middleware"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

import (
	"net/http"
)

func main() {
	apiConfig := initApiConfig()

	mux := http.NewServeMux()

	mux.Handle("/app/", middleware.MetricsInc(apiConfig, http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlers.HandleHealthz)

	mux.HandleFunc("GET /admin/metrics", handlers.HandleAdminMetrics(apiConfig))
	mux.HandleFunc("POST /admin/reset", handlers.HandleAdminReset(apiConfig))

	mux.HandleFunc("POST /api/users", handlers.HandlerCreateUser(apiConfig))

	mux.HandleFunc("POST /api/chirps", handlers.HandleCreateChirp(apiConfig))
	mux.HandleFunc("GET /api/chirps", handlers.HandleListChirps(apiConfig))
	mux.HandleFunc("GET /api/chirps/{chirpID}", handlers.HandleGetChirp(apiConfig))

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func initApiConfig() *config.Api {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = config.Dev
	}

	apiConfig := &config.Api{
		Queries:  database.New(db),
		Platform: platform,
	}
	return apiConfig
}
