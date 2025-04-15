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
	Secret         string
	PolkaKey       string
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

	secret := os.Getenv("SECRET")
	if secret == "" {
		panic("No password signing secret provided")
	}

	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		panic("No polka key provided")
	}

	return &Configuration{
		Queries:  database.New(db),
		Platform: platform,
		Secret:   secret,
		PolkaKey: polkaKey,
	}
}
