package service

import (
	"context"
	"errors"

	"github.com/ericolvr/ewz/internal/domain"
)

type ClientService struct {
	clientRepo domain.ClientRepository
}

func NewClientService(clientRepo domain.ClientRepository) *ClientService {
	return &ClientService{clientRepo: clientRepo}
}

func (s *ClientService) Create(ctx context.Context, client *domain.Client) error {
	if err := client.Validate(); err != nil {
		return err
	}

	existingClient, err := s.clientRepo.FindByEmail(ctx, client.CustomerEmail)
	if err != nil {
		return err
	}
	if existingClient != nil {
		return errors.New("e-mail já cadastrado")
	}

	return s.clientRepo.Create(ctx, client)
}

func (s *ClientService) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	client, err := s.clientRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) UpdateStatus(ctx context.Context, email string, status domain.Status, priority domain.Priority) error {
	return s.clientRepo.UpdateStatus(ctx, email, status, priority)
}
