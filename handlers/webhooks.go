package handlers

import (
	"chirpy/config"
	"chirpy/models"
	"net/http"
)

func HandleWebhook(configuration *config.Configuration) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		webhookRequest, err := decodeRequestPayload[models.WebhookRequest](r)
		if err != nil {
			respondWithError(w, somethingWentWrong, http.StatusBadRequest)
			return
		}

		if webhookRequest.Event != "user.upgraded" {
			printJsonResponse(w, "", http.StatusNoContent)
			return
		}

		user, err := configuration.Queries.GetUser(r.Context(), webhookRequest.Data.UserId)
		if err != nil || len(user.Email) == 0 {
			respondWithError(w, "", http.StatusNotFound)
			return
		}

		if err = configuration.Queries.UpdateUserChirpyRed(r.Context(), user.ID); err != nil {
			respondWithError(w, "", http.StatusNotFound)
			return
		}

		printJsonResponse(w, "", http.StatusNoContent)
	}
}
