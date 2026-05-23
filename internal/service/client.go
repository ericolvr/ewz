package service

import (
	"context"
	"errors"
	"log/slog"

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
		slog.WarnContext(ctx, "validação falhou", "email", client.CustomerEmail, "erro", err.Error())
		return err
	}

	existingClient, err := s.clientRepo.FindByEmail(ctx, client.CustomerEmail)
	if err != nil {
		slog.ErrorContext(ctx, "erro ao buscar cliente", "email", client.CustomerEmail, "erro", err.Error())
		return err
	}
	if existingClient != nil {
		slog.WarnContext(ctx, "e-mail já cadastrado", "email", client.CustomerEmail)
		return errors.New("e-mail já cadastrado")
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		slog.ErrorContext(ctx, "erro ao salvar cliente", "email", client.CustomerEmail, "erro", err.Error())
		return err
	}

	slog.InfoContext(ctx, "cliente criado", "email", client.CustomerEmail, "status", client.Status)

	// Goroutine para nao bloquear a resposta do cliente.
	// Em prod sugeriria usar uma fila SQS para garantir durabilidade,
	// retry automático e observabilidade via DLQ.
	go func() {
		if err := s.pipefyClient.CreateCard(context.Background(), client); err != nil {
			slog.Error("falha ao criar card no pipefy", "email", client.CustomerEmail, "erro", err.Error())
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
