package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"
)

const contentType = "Content-Type"
const textPlainContentType = "text/plain; charset=utf-8"
const applicationJsonContentType = "application/json; charset=utf-8"
const textHtmlContentType = "text/html; charset=utf-8"

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type validationRequest struct {
	Body string `json:"body"`
}

type validationResponse struct {
	Error       *string `json:"error"`
	Valid       *bool   `json:"valid"`
	CleanedBody *string `json:"cleaned_body"`
}

func main() {
	apiConfig := &apiConfig{}
	mux := http.NewServeMux()

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		printResponse(w, "OK", textPlainContentType, http.StatusOK)
	})

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		body := `<html>
		  <body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`

		printResponse(
			w,
			fmt.Sprintf(body, apiConfig.fileServerHits.Load()),
			textHtmlContentType,
			http.StatusOK,
		)
	})

	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		requestBody := validationRequest{}
		err := decoder.Decode(&requestBody)
		if err != nil {
			respondWithError(w, "Something went wrong")
			return
		}

		if len(requestBody.Body) > 140 {
			respondWithError(w, "Chirp is too long")
			return
		}

		text := requestBody.Body
		for _, word := range []string{"kerfuffle ", "sharbert ", "fornax "} {
			text = replaceInsensitive(text, word, "**** ")
		}

		valid := true
		respBody := validationResponse{
			Error:       nil,
			Valid:       &valid,
			CleanedBody: &text,
		}

		printJsonResponse(w, respBody, http.StatusOK)
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileServerHits.Swap(0)
		printResponse(w, "", textPlainContentType, http.StatusOK)
	})

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}

func printJsonResponse(w http.ResponseWriter, respBody validationResponse, statusCode int) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)

		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set(contentType, applicationJsonContentType)
	_, err = w.Write(dat)
	if err != nil {
		panic(err)
	}
}

func respondWithError(w http.ResponseWriter, errorMessage string) {
	printResponse(w, errorMessage, applicationJsonContentType, http.StatusBadRequest)
}

func printResponse(w http.ResponseWriter, response string, contentType string, httpStatus int) {
	w.WriteHeader(httpStatus)

	w.Header().Set(contentType, contentType)

	if len(response) <= 0 {
		return
	}

	if _, err := w.Write([]byte(response)); err != nil {
		panic(err)
	}
}

func replaceInsensitive(input, old, new string) string {
	lowerInput := strings.ToLower(input)
	lowerOld := strings.ToLower(old)
	index := strings.Index(lowerInput, lowerOld)
	if index == -1 {
		return input
	}
	return input[:index] + new + input[index+len(old):]
}
