package pipefy

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ericolvr/ewz/internal/domain"
)

// Mock do envio ao pipefy
// Em prod seria um HTTP real para o endpoint GraphQL do Pipefy.
// Usaria uma fila (SQS) para garantir durabilidade, retry automático e DLQ.

const createCardMutation = `
	mutation CreateCard($pipe_id: ID!, $fields_attributes: [FieldValueInput]!) {
		createCard(input: {
			pipe_id: $pipe_id
			fields_attributes: $fields_attributes
		}) {
			card {
				id
				title
			}
		}
	}
`

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) CreateCard(ctx context.Context, client *domain.Client) error {
	variables := map[string]any{
		"pipe_id": "seu_pipe_id_aqui",
		"fields_attributes": []map[string]any{
			{"field_id": "nome", "field_value": client.CustomerName},
			{"field_id": "email", "field_value": client.CustomerEmail},
			{"field_id": "patrimonio", "field_value": fmt.Sprintf("%.2f", client.AssetValue)},
		},
	}

	payload := map[string]any{
		"query":     createCardMutation,
		"variables": variables,
	}

	out, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	slog.Info("[pipefy] createCard payload", "payload", string(out))
	return nil
}

const updateCardFieldMutation = `
	mutation UpdateCardField($card_id: ID!, $field_id: ID!, $new_value: String!) {
		updateCardField(input: {
			card_id: $card_id
			field_id: $field_id
			new_value: $new_value
		}) {
			card {
				id
				title
			}
		}
	}
`

func (c *Client) UpdateCard(ctx context.Context, cardID string, client *domain.Client) error {
	priority := ""
	if client.Priority != nil {
		priority = string(*client.Priority)
	}

	updates := []map[string]any{
		{"card_id": cardID, "field_id": "status", "new_value": string(client.Status)},
		{"card_id": cardID, "field_id": "prioridade", "new_value": priority},
	}

	for _, vars := range updates {
		payload := map[string]any{
			"query":     updateCardFieldMutation,
			"variables": vars,
		}
		out, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		slog.Info("[pipefy] updateCardField payload", "payload", string(out))
	}
	return nil
}
