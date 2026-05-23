package service

import (
	"context"
	"errors"
	"log"

	"github.com/ericolvr/ewz/internal/domain"
)

type ClientService struct {
	clientRepo   domain.ClientRepository
	pipefyClient domain.PipefyClient
}

func NewClientService(clientRepo domain.ClientRepository, pipefyClient domain.PipefyClient) *ClientService {
	return &ClientService{clientRepo: clientRepo, pipefyClient: pipefyClient}
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

	if err := s.clientRepo.Create(ctx, client); err != nil {
		return err
	}

	// Goroutine para nao bloquear a resposta do cliente
	// Em prod sugeriria usar uma fila SQS para garantir durabilidade,
	// retry automático e observabilidade via DLQ.

	go func() {
		if err := s.pipefyClient.CreateCard(context.Background(), client); err != nil {
			log.Printf("[pipefy] falha ao criar card para %s: %v", client.CustomerEmail, err)
		}
	}()

	return nil
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
