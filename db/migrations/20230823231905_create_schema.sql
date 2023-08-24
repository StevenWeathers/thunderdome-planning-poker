-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS public.schema_migrations;
CREATE SCHEMA IF NOT EXISTS thunderdome;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA thunderdome;
-- +goose StatementEnd
