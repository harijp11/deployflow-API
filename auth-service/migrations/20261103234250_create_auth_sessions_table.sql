-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY,
    credential_id UUID NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    family_id UUID NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    last_used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_auth_sessions_expiry CHECK (expires_at > created_at)
);

CREATE INDEX idx_auth_sessions_credential_id ON auth_sessions(credential_id);
CREATE INDEX idx_auth_sessions_family_id ON auth_sessions(family_id);
CREATE INDEX idx_auth_sessions_expires_at ON auth_sessions(expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth_sessions;
-- +goose StatementEnd