package config

import (
	"chirpy/internal/database"
	"database/sql"
	"github.com/joho/godotenv"
	"os"
	"sync/atomic"
)

type Configuration struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
}

const (
	Dev  string = "DEV"
	Test string = "TEST"
)

func Init() *Configuration {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		platform = Dev
	}

	return &Configuration{
		Queries:  database.New(db),
		Platform: platform,
	}
}
