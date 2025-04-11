package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiConfig := &apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		printResponse(w, "OK", "text/plain; charset=utf-8")
	})
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		body := `<html>
		  <body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`

		printResponse(w, fmt.Sprintf(body, apiConfig.fileServerHits.Load()), "text/html; charset=utf-8")
	})
	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileServerHits.Swap(0)
		printResponse(w, "", "text/plain; charset=utf-8")
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

func printResponse(w http.ResponseWriter, response string, contentType string) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", contentType)

	if len(response) <= 0 {
		return
	}

	if _, err := w.Write([]byte(response)); err != nil {
		panic(err)
	}
}
