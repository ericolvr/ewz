package service

import (
	"context"
	"testing"

	"github.com/ericolvr/ewz/internal/domain"
)

// --- mocks ---

type mockWebhookRepo struct {
	events map[string]*domain.WebhookEvent
}

func newMockWebhookRepo() *mockWebhookRepo {
	return &mockWebhookRepo{events: make(map[string]*domain.WebhookEvent)}
}

func (m *mockWebhookRepo) FindByEventID(ctx context.Context, eventID string) (*domain.WebhookEvent, error) {
	return m.events[eventID], nil
}

func (m *mockWebhookRepo) Create(ctx context.Context, event *domain.WebhookEvent) error {
	m.events[event.EventID] = event
	return nil
}

func newWebhookSvc(assetValue float64) (*WebhookService, *mockClientRepo) {
	clientRepo := newMockClientRepo()
	clientRepo.clients["joao@example.com"] = &domain.Client{
		CustomerName:  "João Silva",
		CustomerEmail: "joao@example.com",
		RequestType:   "Atualização Cadastral",
		AssetValue:    assetValue,
		Status:        domain.StatusAguardandoAnalise,
	}
	webhookRepo := newMockWebhookRepo()
	svc := NewWebhookService(clientRepo, webhookRepo, &mockPipefyClient{})
	return svc, clientRepo
}

// --- testes ---

func TestWebhook_PrioridadeAlta(t *testing.T) {
	svc, clientRepo := newWebhookSvc(250000)

	if err := svc.Process(context.Background(), "evt_1", "card_1", "joao@example.com"); err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}

	client := clientRepo.clients["joao@example.com"]
	if client.Priority == nil || *client.Priority != domain.PriorityAlta {
		t.Errorf("esperava prioridade_alta, got %v", client.Priority)
	}
	if client.Status != domain.StatusProcessado {
		t.Errorf("esperava status Processado, got %q", client.Status)
	}
}

func TestWebhook_PrioridadeNormal(t *testing.T) {
	svc, clientRepo := newWebhookSvc(100000)

	if err := svc.Process(context.Background(), "evt_2", "card_2", "joao@example.com"); err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}

	client := clientRepo.clients["joao@example.com"]
	if client.Priority == nil || *client.Priority != domain.PriorityNormal {
		t.Errorf("esperava prioridade_normal, got %v", client.Priority)
	}
}

func TestWebhook_EventIDDuplicado(t *testing.T) {
	svc, _ := newWebhookSvc(250000)

	if err := svc.Process(context.Background(), "evt_dup", "card_1", "joao@example.com"); err != nil {
		t.Fatalf("primeiro processamento falhou: %v", err)
	}

	if err := svc.Process(context.Background(), "evt_dup", "card_1", "joao@example.com"); err == nil {
		t.Fatal("esperava erro de evento duplicado, got nil")
	}
}
