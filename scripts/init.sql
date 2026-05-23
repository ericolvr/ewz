CREATE TYPE client_status AS ENUM (
    'Aguardando Análise',
    'Processado'
);

CREATE TYPE client_priority AS ENUM (
    'prioridade_alta',
    'prioridade_normal'
);

CREATE TABLE IF NOT EXISTS clients (
    id               SERIAL PRIMARY KEY,
    cliente_nome     VARCHAR(255)        NOT NULL,
    cliente_email    VARCHAR(255)        NOT NULL UNIQUE,
    tipo_solicitacao VARCHAR(255)        NOT NULL,
    valor_patrimonio NUMERIC(15,2)       NOT NULL,
    status           client_status       NOT NULL DEFAULT 'Aguardando Análise',
    prioridade       client_priority     NULL,
    created_at       TIMESTAMP           NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP           NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_events (
    id            SERIAL PRIMARY KEY,
    event_id      VARCHAR(255) NOT NULL UNIQUE,
    card_id       VARCHAR(255) NOT NULL,
    cliente_email VARCHAR(255) NOT NULL,
    processed_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);
