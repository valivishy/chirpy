package handlers

import (
	"chirpy/config"
	"fmt"
	"net/http"
)

func HandleAdminMetrics(apiConfig *config.Api) func(w http.ResponseWriter, r *http.Request) {
	body := `<html>
		  <body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		  </body>
		</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		printResponse(
			w,
			fmt.Sprintf(body, apiConfig.FileServerHits.Load()),
			textHtmlContentType,
			http.StatusOK,
		)
	}
}

func HandleAdminReset(apiConfig *config.Api) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		apiConfig.FileServerHits.Swap(0)
		printResponse(w, "", textPlainContentType, http.StatusOK)
	}
}
