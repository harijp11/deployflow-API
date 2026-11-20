-- +goose Up
-- +goose StatementBegin

ALTER TABLE auth_sessions
    ADD COLUMN device_name VARCHAR(100),
    ADD COLUMN device_type VARCHAR(50),
    ADD COLUMN browser VARCHAR(100),
    ADD COLUMN operating_system VARCHAR(100),
    ADD COLUMN ip_address VARCHAR(45),
    ADD COLUMN user_agent TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE auth_sessions
    DROP COLUMN device_name,
    DROP COLUMN device_type,
    DROP COLUMN browser,
    DROP COLUMN operating_system,
    DROP COLUMN ip_address,
    DROP COLUMN user_agent;

-- +goose StatementEnd