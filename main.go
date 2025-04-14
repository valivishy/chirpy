package main

import (
	"chirpy/config"
	"chirpy/internal/database"
	"chirpy/router"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
)

import (
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: router.New(initApiConfig()),
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
