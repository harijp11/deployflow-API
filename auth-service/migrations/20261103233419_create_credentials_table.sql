-- +goose Up
-- +goose StatementBegin
CREATE TABLE credentials (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    email VARCHAR(320) NOT NULL,
    password_hash TEXT NOT NULL,
    email_verified_at TIMESTAMP,
    password_changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_credentials_email_not_empty CHECK (LENGTH(TRIM(email)) > 0),
    CONSTRAINT chk_credentials_password_hash_not_empty CHECK (LENGTH(password_hash) > 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE credentials;
-- +goose StatementEnd