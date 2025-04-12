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

	mux.HandleFunc("POST /api/validate_chirp", handlers.HandleValidateChirp)

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

	apiConfig := &config.Api{
		Queries: database.New(db),
	}
	return apiConfig
}
