package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ericolvr/ewz/internal/domain"
)

type ClientRepository struct {
	db *sql.DB
}

func NewClientRepository(db *sql.DB) domain.ClientRepository {
	return &ClientRepository{db: db}
}

func (c *ClientRepository) Create(ctx context.Context, client *domain.Client) error {
	query := `
		INSERT INTO clients (
			cliente_nome,
			cliente_email,
			tipo_solicitacao,
			valor_patrimonio,
			status,
			prioridade
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`
	return c.db.QueryRowContext(
		ctx,
		query,
		client.CustomerName,
		client.CustomerEmail,
		client.RequestType,
		client.AssetValue,
		client.Status,
		client.Priority,
	).Scan(&client.ID)
}

func (c *ClientRepository) FindByEmail(ctx context.Context, email string) (*domain.Client, error) {
	query := `
		SELECT
			id,
			cliente_nome,
			cliente_email,
			tipo_solicitacao,
			valor_patrimonio,
			status,
			prioridade
		FROM clients
		WHERE cliente_email = $1
	`
	row := c.db.QueryRowContext(ctx, query, email)

	var client domain.Client
	var priority sql.NullString
	err := row.Scan(
		&client.ID,
		&client.CustomerName,
		&client.CustomerEmail,
		&client.RequestType,
		&client.AssetValue,
		&client.Status,
		&priority,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if priority.Valid {
		p := domain.Priority(priority.String)
		client.Priority = &p
	}
	return &client, nil
}

func (c *ClientRepository) UpdateStatus(ctx context.Context, email string, status domain.Status, priority domain.Priority) error {
	query := `
		UPDATE clients
		SET status = $1, prioridade = $2, updated_at = NOW()
		WHERE cliente_email = $3
	`
	_, err := c.db.ExecContext(ctx, query, status, priority, email)
	return err
}
