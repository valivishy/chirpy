package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

const contentTypeHeaderName = "Content-Type"
const applicationJsonContentType = "application/json; charset=utf-8"
const textPlainContentType = "text/plain; charset=utf-8"
const textHtmlContentType = "text/html; charset=utf-8"

func respondWithError(w http.ResponseWriter, errorMessage string, status int) {
	printResponse(w, errorMessage, applicationJsonContentType, status)
}

func printJsonResponse(w http.ResponseWriter, respBody any, statusCode int) {
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)

		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set(contentTypeHeaderName, applicationJsonContentType)
	_, err = w.Write(dat)
	if err != nil {
		panic(err)
	}
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
