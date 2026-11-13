-- +goose Up
-- +goose StatementBegin
CREATE TABLE email_verification_tokens (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    credential_id INT NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE email_verification_tokens;
-- +goose StatementEnd