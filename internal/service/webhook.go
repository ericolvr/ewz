package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ericolvr/ewz/internal/domain"
)

type WebhookService struct {
	clientRepo   domain.ClientRepository
	webhookRepo  domain.WebhookEventRepository
	pipefyClient domain.PipefyClient
}

func NewWebhookService(
	clientRepo domain.ClientRepository,
	webhookRepo domain.WebhookEventRepository,
	pipefyClient domain.PipefyClient,
) *WebhookService {
	return &WebhookService{
		clientRepo:   clientRepo,
		webhookRepo:  webhookRepo,
		pipefyClient: pipefyClient,
	}
}

func (s *WebhookService) Process(ctx context.Context, eventID, cardID, clientEmail string) error {
	slog.InfoContext(ctx, "webhook recebido", "event_id", eventID, "card_id", cardID, "email", clientEmail)

	existing, err := s.webhookRepo.FindByEventID(ctx, eventID)
	if err != nil {
		slog.ErrorContext(ctx, "erro ao verificar idempotência", "event_id", eventID, "erro", err.Error())
		return err
	}
	if existing != nil {
		slog.WarnContext(ctx, "evento duplicado bloqueado", "event_id", eventID)
		return errors.New("evento já processado")
	}

	client, err := s.clientRepo.FindByEmail(ctx, clientEmail)
	if err != nil {
		slog.ErrorContext(ctx, "erro ao buscar cliente", "email", clientEmail, "erro", err.Error())
		return err
	}
	if client == nil {
		slog.WarnContext(ctx, "cliente não encontrado", "email", clientEmail)
		return errors.New("cliente não encontrado")
	}

	priority := domain.PriorityNormal
	if client.AssetValue >= 200000 {
		priority = domain.PriorityAlta
	}

	client.Status = domain.StatusProcessado
	client.Priority = &priority

	if err := s.clientRepo.UpdateStatus(ctx, clientEmail, client.Status, priority); err != nil {
		slog.ErrorContext(ctx, "erro ao atualizar status", "email", clientEmail, "erro", err.Error())
		return err
	}

	if err := s.pipefyClient.UpdateCard(ctx, cardID, client); err != nil {
		slog.ErrorContext(ctx, "erro ao atualizar card no pipefy", "card_id", cardID, "erro", err.Error())
		return err
	}

	event := &domain.WebhookEvent{
		EventID:     eventID,
		CardID:      cardID,
		ClientEmail: clientEmail,
	}
	if err := s.webhookRepo.Create(ctx, event); err != nil {
		slog.ErrorContext(ctx, "erro ao salvar evento", "event_id", eventID, "erro", err.Error())
		return err
	}

	slog.InfoContext(ctx, "webhook processado", "event_id", eventID, "email", clientEmail, "prioridade", priority)
	return nil
}
