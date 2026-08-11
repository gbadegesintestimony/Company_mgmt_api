DROP INDEX IF EXISTS idx_sessions_refresh_hash;

ALTER TABLE sessions
    ADD CONSTRAINT uq_sessions_refresh_token_hash UNIQUE (refresh_token_hash);
