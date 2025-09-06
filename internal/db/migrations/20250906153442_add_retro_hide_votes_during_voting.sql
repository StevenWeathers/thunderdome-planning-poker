-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.retro ADD COLUMN hide_votes_during_voting BOOLEAN NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.retro DROP COLUMN IF EXISTS hide_votes_during_voting;
-- +goose StatementEnd
