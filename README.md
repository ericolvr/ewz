# EWZ — Client Management & Pipefy Integration

API em Go para gerenciamento de clientes e integração com o Pipefy via GraphQL.

---

## Requisitos

- Go 1.21+
- Docker e Docker Compose

---

## Execução local

```bash
# 1. Instalar dependências
make install

# 2. Subir o banco de dados (cria as tabelas automaticamente)
make db-start

# 3. Rodar o servidor
make run
```

A API estará disponível em `http://localhost:8080`.

---

## Testes

```bash
make test
```

---

## Endpoints

### POST /api/v1/clientes

Cria um novo cliente e simula a criação de um card no Pipefy.

```bash
curl -X POST http://localhost:8080/api/v1/clientes \
  -H "Content-Type: application/json" \
  -d '{
    "cliente_nome": "João Silva",
    "cliente_email": "joao.silva@example.com",
    "tipo_solicitacao": "Atualização Cadastral",
    "valor_patrimonio": 250000
  }'
```

**Resposta:**
```json
{
  "id": 1,
  "cliente_nome": "João Silva",
  "cliente_email": "joao.silva@example.com",
  "tipo_solicitacao": "Atualização Cadastral",
  "valor_patrimonio": 250000,
  "status": "Aguardando Análise"
}
```

---

### POST /api/v1/webhooks/pipefy/card-updated

Simula o recebimento de um webhook do Pipefy quando um operacional atualiza um card.
Define a prioridade do cliente com base no patrimônio e atualiza o status para `Processado`.

```bash
curl -X POST http://localhost:8080/api/v1/webhooks/pipefy/card-updated \
  -H "Content-Type: application/json" \
  -d '{
    "event_id": "evt_123",
    "card_id": "card_456",
    "cliente_email": "joao.silva@example.com",
    "timestamp": "2026-05-18T12:00:00Z"
  }'
```

**Resposta:**
```json
{
  "message": "evento processado com sucesso"
}
```

**Regra de prioridade:**
- `valor_patrimonio >= 200.000` → `prioridade_alta`
- `valor_patrimonio < 200.000` → `prioridade_normal`

**Idempotência:** requisições com o mesmo `event_id` retornam `409 Conflict`.

---

## Comandos disponíveis

| Comando | Descrição |
|---|---|
| `make install` | Instala dependências |
| `make db-start` | Sobe o banco de dados |
| `make db-stop` | Para o banco de dados |
| `make db-reset` | Recria o banco do zero |
| `make run` | Roda o servidor |
| `make test` | Executa os testes |
| `make build` | Gera o binário |

---

## Visão de Produção (AWS)

Em produção, a arquitetura escalaria da seguinte forma:

**Criação de cliente (`POST /clientes`)**

A API rodaria em **AWS Lambda** exposta via **API Gateway**. Ao criar um cliente, após persistir no **RDS (PostgreSQL)**, a Lambda publicaria uma mensagem no **SQS** em vez de chamar o Pipefy diretamente. Uma segunda Lambda consumiria a fila e faria a chamada GraphQL ao Pipefy com retry automático. Falhas que esgotassem as tentativas iriam para uma **Dead Letter Queue (DLQ)** para observabilidade e reprocessamento manual.

**Webhook (`POST /webhooks/pipefy/card-updated`)**

O endpoint de webhook também rodaria via API Gateway + Lambda. A idempotência seria garantida por uma tabela no **DynamoDB** (chave: `event_id`), que oferece latência baixa para esse tipo de lookup pontual. O cliente seria buscado no RDS e atualizado após o cálculo de prioridade.

**Benefícios:**
- Escalabilidade automática via Lambda
- Desacoplamento da integração com Pipefy via SQS
- Garantia de entrega com retry e DLQ
- Idempotência robusta via DynamoDB
