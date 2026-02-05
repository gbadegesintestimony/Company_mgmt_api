CREATE INDEX idx_users_company_email
ON users (company_id, email);

CREATE INDEX idx_users_active
ON users (company_id)
WHERE is_active = true;

CREATE INDEX idx_otps_active
ON otps (user_id, purpose, expires_at)
WHERE consumed_at IS NULL;
