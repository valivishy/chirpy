package models

import "github.com/google/uuid"

type WebhookRequest struct {
	Event string      `json:"event"`
	Data  WebhookData `json:"data"`
}

type WebhookData struct {
	UserId uuid.UUID `json:"user_id"`
}
