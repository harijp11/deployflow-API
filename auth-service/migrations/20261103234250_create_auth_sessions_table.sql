-- +goose Up
-- +goose StatementBegin
CREATE TABLE auth_sessions (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    credential_id INT NOT NULL REFERENCES credentials(id) ON DELETE CASCADE,
    user_id INT NOT NULL,
    token_hash VARCHAR(255) NOT NULL UNIQUE,
    family_id INT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP,
    last_used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE auth_sessions;
-- +goose StatementEnd