package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	bearerToken := headers.Get("Authorization")
	if len(bearerToken) == 0 {
		return "", fmt.Errorf("no bearer token found")
	}

	return strings.TrimSpace(strings.Replace(bearerToken, "ApiKey ", "", 1)), nil
}
