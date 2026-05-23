# EWZ — Client Management & Pipefy Integration

API em Go para gerenciamento de clientes e integração com o Pipefy via GraphQL.

---

## Requisitos

- Go 1.21+
- Docker e Docker Compose

---

## Execução local
Todos os comandos estão no makfile e podem ser vistos com make help
```bash
# 1. Instalar dependências
make install

# 2. Subir o banco de dados (cria as tabelas automaticamente)make db-start

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

Em produção, o sistema escalaria com a seguinte arquitetura:

**Aplicação**

Cada endpoint rodaria como uma função **Lambda** exposta via **API Gateway**. Essa abordagem elimina o gerenciamento de servidores, escala por requisição e cobra apenas pelo uso. O Gin precisaria de um adapter (`aws-lambda-go`) para rodar nesse ambiente. Como alternativa para times que preferem servidor dedicado, a API poderia rodar em EC2 dentro de um Auto Scaling Group (ASG) com um Application Load Balancer (ALB) na frente.

**Banco de dados**

O **RDS PostgreSQL** em configuração Multi-AZ mantém uma réplica em standby em outra zona de disponibilidade. Em caso de falha da instância primária, o RDS promove a réplica automaticamente, sem perda de dados. A conexão é configurada via variável de ambiente, bastando apontar para o endpoint do RDS.

**Integração com Pipefy**

Ao criar um cliente, a Lambda publicaria uma mensagem no **SQS** em vez de chamar o Pipefy diretamente — desacoplando a integração e garantindo durabilidade. Uma segunda Lambda consumiria a fila e executaria a mutation GraphQL com retry automático. Falhas que esgotassem as tentativas iriam para uma **Dead Letter Queue (DLQ)**, permitindo reprocessamento manual e observabilidade.

**Webhook**

A Lambda do webhook recebe o evento do Pipefy, verifica idempotência via **DynamoDB** (chave: `event_id`) — oferecendo latência baixa para esse tipo de lookup pontual — calcula a prioridade do cliente e persiste o resultado no RDS. A chamada de atualização ao Pipefy (`updateCardField`) é feita de forma síncrona; em produção, poderia ser extraída para uma fila SQS caso o tempo de resposta se tornasse crítico.
