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
		) VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := c.db.ExecContext(
		ctx,
		query,
		client.CustomerName,
		client.CustomerEmail,
		client.RequestType,
		client.AssetValue,
		client.Status,
		client.Priority,
	)
	return err
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
		WHERE cliente_email = ?
	`
	row := c.db.QueryRowContext(ctx, query, email)
	var client domain.Client
	err := row.Scan(
		&client.ID,
		&client.CustomerName,
		&client.CustomerEmail,
		&client.RequestType,
		&client.AssetValue,
		&client.Status,
		&client.Priority,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

func (c *ClientRepository) UpdateStatus(ctx context.Context, email, status, priority string) error {
	query := `
		UPDATE clients
		SET status = ?, prioridade = ?
		WHERE cliente_email = ?
	`
	_, err := c.db.ExecContext(ctx, query, status, priority, email)
	return err
}
