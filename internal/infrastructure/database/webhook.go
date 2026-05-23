package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/ewz/internal/domain"
)

type WebhookEventRepository struct {
	db *sql.DB
}

func NewWebhookEventRepository(db *sql.DB) domain.WebhookEventRepository {
	return &WebhookEventRepository{db: db}
}

func (r *WebhookEventRepository) FindByEventID(ctx context.Context, eventID string) (*domain.WebhookEvent, error) {
	query := `SELECT id, event_id, card_id, cliente_email, processed_at FROM webhook_events WHERE event_id = $1`
	row := r.db.QueryRowContext(ctx, query, eventID)

	var event domain.WebhookEvent
	err := row.Scan(&event.ID, &event.EventID, &event.CardID, &event.ClientEmail, &event.ProcessedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (r *WebhookEventRepository) Create(ctx context.Context, event *domain.WebhookEvent) error {
	query := `INSERT INTO webhook_events (event_id, card_id, cliente_email) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, event.EventID, event.CardID, event.ClientEmail)
	return err
}
