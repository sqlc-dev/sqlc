-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA auth;
CREATE EXTENSION "pgcrypto";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA auth;
-- +goose StatementEnd
