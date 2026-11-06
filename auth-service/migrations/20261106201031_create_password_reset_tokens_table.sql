-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_reset_tokens (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,   
    credential_id INT NOT NULL,
    token_hash VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE password_reset_tokens;
-- +goose StatementEnd