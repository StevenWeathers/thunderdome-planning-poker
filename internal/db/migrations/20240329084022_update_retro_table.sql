-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.retro ADD COLUMN ready_users jsonb DEFAULT '[]'::jsonb;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.retro DROP COLUMN ready_users;
-- +goose StatementEnd
