package tests

import (
	"chirpy/config"
	"chirpy/models"
	"github.com/google/uuid"
	"net/http"
	"testing"
)

func TestHandleWebhook_WrongStatus(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	request := models.WebhookRequest{
		Event: "user.downgraded",
		Data: models.WebhookData{
			UserId: uuid.New(),
		},
	}

	execPost(t, ts, "/api/polka/webhooks", stringify(request), config.Init().PolkaKey, "ApiKey", http.StatusNoContent, &struct{}{})

	email := "test-webhook-wrong-status@example.com"
	password := "superPassword123"
	userId := createUser(t, ts, email, password)
	request.Data.UserId = userId

	execPost(t, ts, "/api/polka/webhooks", stringify(request), config.Init().PolkaKey, "ApiKey", http.StatusNoContent, &struct{}{})
}

func TestHandleWebhook_RightStatusWrongUser(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	request := models.WebhookRequest{
		Event: "user.upgraded",
		Data: models.WebhookData{
			UserId: uuid.New(),
		},
	}

	execPost(t, ts, "/api/polka/webhooks", stringify(request), config.Init().PolkaKey, "ApiKey", http.StatusNotFound, &struct{}{})
}

func TestHandleWebhook_RightStatusRightUser(t *testing.T) {
	ts := Start(t)
	defer closer(t)(ts.Server)

	email := "test-webhook-right-user-right-status@example.com"
	password := "superPassword123"
	userId := createUser(t, ts, email, password)

	request := models.WebhookRequest{
		Event: "user.upgraded",
		Data: models.WebhookData{
			UserId: userId,
		},
	}

	execPost(t, ts, "/api/polka/webhooks", stringify(request), config.Init().PolkaKey, "ApiKey", http.StatusNoContent, &struct{}{})
}
