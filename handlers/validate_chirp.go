package handlers

import (
	"encoding/json"
	"net/http"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	requestBody := ValidationRequest{}
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
	respBody := ValidationResponse{
		Error:       nil,
		Valid:       &valid,
		CleanedBody: &text,
	}

	printJsonResponse(w, respBody, http.StatusOK)
}
