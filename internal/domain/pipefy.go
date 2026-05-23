package domain

import "context"

type PipefyClient interface {
	CreateCard(ctx context.Context, client *Client) error
	UpdateCard(ctx context.Context, cardID string, client *Client) error
}
