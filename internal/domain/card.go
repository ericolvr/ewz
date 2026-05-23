package domain

import (
	"context"
	"time"
)

type WebhookEvent struct {
	ID          int64     `json:"id" db:"id"`
	EventID     string    `json:"event_id" db:"event_id"`
	CardID      string    `json:"card_id" db:"card_id"`
	ClientEmail string    `json:"cliente_email" db:"cliente_email"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

type WebhookEventRepository interface {
	FindByEventID(ctx context.Context, eventID string) (*WebhookEvent, error)
	Create(ctx context.Context, event *WebhookEvent) error
}
