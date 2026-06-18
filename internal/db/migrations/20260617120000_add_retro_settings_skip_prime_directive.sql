-- +goose Up
-- +goose StatementBegin

ALTER TABLE thunderdome.retro_settings
ADD COLUMN skip_prime_directive BOOLEAN NOT NULL DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE thunderdome.retro_settings
DROP COLUMN skip_prime_directive;

-- +goose StatementEnd