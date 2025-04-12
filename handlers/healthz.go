package handlers

import "net/http"

func HandleHealthz(w http.ResponseWriter, _ *http.Request) {
	printResponse(w, "OK", textPlainContentType, http.StatusOK)
}
