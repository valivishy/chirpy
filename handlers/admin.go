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
		if apiConfig.Platform != config.Dev {
			printResponse(w, "", textPlainContentType, http.StatusForbidden)
			return
		}

		apiConfig.FileServerHits.Swap(0)
		if err := apiConfig.Queries.DeleteUsers(r.Context()); err != nil {
			printResponse(w, "", textHtmlContentType, http.StatusInternalServerError)
			return
		}

		if err := apiConfig.Queries.DeleteChirps(r.Context()); err != nil {
			printResponse(w, "", textHtmlContentType, http.StatusInternalServerError)
			return
		}

		printResponse(w, "", textPlainContentType, http.StatusOK)
	}
}
