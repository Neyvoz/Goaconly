CREATE TABLE IF NOT EXISTS check_logs (
    id               BIGSERIAL PRIMARY KEY,
    target_id        BIGINT NOT NULL REFERENCES targets(id) ON DELETE CASCADE,
    checked_at       TIMESTAMPTZ NOT NULL,
    status_code      INT NOT NULL DEFAULT 0,
    response_time_ms BIGINT NOT NULL DEFAULT 0,
    is_up            BOOLEAN NOT NULL,
    ssl_expires_at   TIMESTAMPTZ,
    error_message    TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_checl_log_target_id_checked_at
    ON check_logs (target_id, checked_at DESC);