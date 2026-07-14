-- Расширение для генерации UUID (пригодится в будущем)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Тарифные планы
CREATE TABLE IF NOT EXISTS plans (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,  -- 'free', 'starter', 'pro'
    max_targets INTEGER NOT NULL DEFAULT 5,
    price_cents INTEGER NOT NULL DEFAULT 0
);

-- Пользователи
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    plan_id       INTEGER NOT NULL REFERENCES plans(id) DEFAULT 1,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    company_name  VARCHAR(255) NOT NULL DEFAULT ''
);

-- Цели мониторинга
CREATE TABLE IF NOT EXISTS targets (
    id              BIGSERIAL PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url             VARCHAR(2048) NOT NULL,
    check_interval  SMALLINT NOT NULL DEFAULT 5 CHECK (check_interval BETWEEN 1 AND 1440),
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Составной уникальный индекс: один юзер не может добавить один URL дважды
    CONSTRAINT uq_user_url UNIQUE (user_id, url)
);

-- Логи проверок (высоконагруженная таблица — будет расти быстро)
CREATE TABLE IF NOT EXISTS check_logs (
    id             BIGSERIAL PRIMARY KEY,
    target_id      BIGINT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    checked_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status         VARCHAR(20) NOT NULL CHECK (status IN ('up', 'down', 'pending_down')),
    status_code    SMALLINT,
    response_time  INTEGER,  -- в миллисекундах
    error_message  TEXT,
    ssl_expires_at TIMESTAMPTZ
);

-- Рефреш токен
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Поиск по хэшу при рефреше
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens (token_hash);

-- Массовый revoke всех токенов юзера
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens (user_id) WHERE revoked = FALSE;

-- Индексы производительности
-- Запрос "дай последние N логов для target X" — самый частый
CREATE INDEX idx_check_logs_target_time ON check_logs (target_id, checked_at DESC);

-- Запрос "все активные цели конкретного юзера"
CREATE INDEX idx_targets_user_active ON targets (user_id, is_active);

-- Данные для seed
INSERT INTO plans (name, max_targets, price_cents) VALUES
    ('free',    3,   0),
    ('starter', 20,  1900),
    ('pro',     100, 4900);