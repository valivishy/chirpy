package handlers

type ValidationRequest struct {
	Body string `json:"body"`
}

type ValidationResponse struct {
	Error       *string `json:"error"`
	Valid       *bool   `json:"valid"`
	CleanedBody *string `json:"cleaned_body"`
}
