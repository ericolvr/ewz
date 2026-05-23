package dto

type WebhookRequest struct {
	EventID     string `json:"event_id"      binding:"required"`
	CardID      string `json:"card_id"       binding:"required"`
	ClientEmail string `json:"cliente_email" binding:"required,email"`
	Timestamp   string `json:"timestamp"     binding:"required"`
}
