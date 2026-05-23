package domain

import (
	"context"
	"errors"
	"regexp"
)

var mailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Client struct {
	ID            int64   `json:"id" db:"id"`
	CustomerName  string  `json:"cliente_nome" db:"cliente_nome"`
	CustomerEmail string  `json:"cliente_email" db:"cliente_email"`
	RequestType   string  `json:"tipo_solicitacao" db:"tipo_solicitacao"`
	AssetValue    float64 `json:"valor_patrimonio" db:"valor_patrimonio"`
	Status        string  `json:"status" db:"status"`
	Priority      string  `json:"prioridade" db:"prioridade"`
}

type ClientRepository interface {
	Create(ctx context.Context, client *Client) error
	FindByEmail(ctx context.Context, email string) (*Client, error)
	UpdateStatus(ctx context.Context, email, status, priority string) error
}

func (c *Client) Validate() error {
	if c.CustomerName == "" {
		return errors.New("nome do cliente é obrigatório")
	}
	if c.CustomerEmail == "" {
		return errors.New("e-mail do cliente é obrigatório")
	}
	if !mailRegex.MatchString(c.CustomerEmail) {
		return errors.New("e-mail do cliente é inválido")
	}
	if c.RequestType == "" {
		return errors.New("tipo de solicitação é obrigatório")
	}
	if c.AssetValue == 0 {
		return errors.New("valor do patrimônio é obrigatório")
	}
	return nil
}
