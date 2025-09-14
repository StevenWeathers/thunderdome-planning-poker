-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.poker
ADD COLUMN IF NOT EXISTS end_time TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS end_reason VARCHAR(16); -- Possible values: 'completed', 'abandoned', 'cancelled', etc.
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.poker
DROP COLUMN IF EXISTS end_time,
DROP COLUMN IF EXISTS end_reason;
-- +goose StatementEnd
