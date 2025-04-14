package main

import (
	"chirpy/config"
	"chirpy/router"
	_ "github.com/lib/pq"
)

import (
	"net/http"
)

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: router.New(config.Init()),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
