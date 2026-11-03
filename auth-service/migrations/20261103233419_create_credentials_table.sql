-- +goose Up
-- +goose StatementBegin
CREATE TABLE credentials (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL,
    password_hash TEXT NOT NULL,
    email_verified_at TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_credentials_email_not_empty CHECK (LENGTH(TRIM(email)) > 0),
    CONSTRAINT chk_credentials_password_hash_not_empty CHECK (LENGTH(password_hash) > 0)
);

CREATE UNIQUE INDEX idx_credentials_email_unique ON credentials (LOWER(email));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE credentials;
-- +goose StatementEnd