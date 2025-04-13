package config

import (
	"chirpy/internal/database"
	"sync/atomic"
)

type Api struct {
	FileServerHits atomic.Int32
	Queries        *database.Queries
	Platform       string
}

const (
	Dev string = "DEV"
)
