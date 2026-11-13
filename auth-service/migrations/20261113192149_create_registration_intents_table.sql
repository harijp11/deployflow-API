-- +goose Up
-- +goose StatementBegin
CREATE TABLE registration_intents (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    company_name VARCHAR(255) NULL,
    phone_number VARCHAR(20) NULL,
    verification_token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    email_send_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    email_send_attempts INT NOT NULL DEFAULT 0,
    last_email_sent_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE registration_intents;
-- +goose StatementEnd