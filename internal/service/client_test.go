package service

import (
	"context"
	"testing"

	"github.com/ericolvr/ewz/internal/domain"
)

// --- mocks ---

type mockClientRepo struct {
	clients map[string]*domain.Client
}

func newMockClientRepo() *mockClientRepo {
	return &mockClientRepo{clients: make(map[string]*domain.Client)}
}

func (m *mockClientRepo) Create(ctx context.Context, client *domain.Client) error {
	m.clients[client.CustomerEmail] = client
	return nil
}

func (m *mockClientRepo) FindByEmail(ctx context.Context, email string) (*domain.Client, error) {
	return m.clients[email], nil
}

func (m *mockClientRepo) UpdateStatus(ctx context.Context, email string, status domain.Status, priority domain.Priority) error {
	if c, ok := m.clients[email]; ok {
		c.Status = status
		c.Priority = &priority
	}
	return nil
}

type mockPipefyClient struct{}

func (m *mockPipefyClient) CreateCard(ctx context.Context, client *domain.Client) error {
	return nil
}

func (m *mockPipefyClient) UpdateCard(ctx context.Context, cardID string, client *domain.Client) error {
	return nil
}

// --- testes ---

func TestCreateClient_ValidPayload(t *testing.T) {
	repo := newMockClientRepo()
	svc := NewClientService(repo, &mockPipefyClient{})

	client := &domain.Client{
		CustomerName:  "João Silva",
		CustomerEmail: "joao@example.com",
		RequestType:   "Atualização Cadastral",
		AssetValue:    250000,
		Status:        domain.StatusAguardandoAnalise,
	}

	if err := svc.Create(context.Background(), client); err != nil {
		t.Fatalf("esperava nil, got: %v", err)
	}

	saved, _ := repo.FindByEmail(context.Background(), "joao@example.com")
	if saved == nil {
		t.Fatal("cliente não foi salvo no banco")
	}
	if saved.Status != domain.StatusAguardandoAnalise {
		t.Errorf("status esperado %q, got %q", domain.StatusAguardandoAnalise, saved.Status)
	}
}
