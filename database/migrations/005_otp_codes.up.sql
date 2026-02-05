-- Enable UUID generation if not already enabled
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE otps (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    purpose TEXT NOT NULL,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    consumed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_otps_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Fast lookup for OTP validation
CREATE INDEX idx_otps_user_purpose
    ON otps (user_id, purpose);

-- Optional: expire lookup
CREATE INDEX idx_otps_expires_at
    ON otps (expires_at);
