ALTER TABLE sessions
    DROP CONSTRAINT IF EXISTS uq_sessions_refresh_token_hash;

CREATE INDEX IF NOT EXISTS idx_sessions_refresh_hash ON sessions(refresh_token_hash);
