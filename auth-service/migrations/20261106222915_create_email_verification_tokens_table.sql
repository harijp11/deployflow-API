
-- +goose Up

-- +goose StatementBegin
CREATE TABLE email_verification_tokens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    credential_id BIGINT NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS email_verification_tokens;
-- +goose StatementEnd
