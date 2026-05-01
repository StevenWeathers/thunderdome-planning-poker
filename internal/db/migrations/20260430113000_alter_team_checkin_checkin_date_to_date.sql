-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.team_checkin
    ALTER COLUMN checkin_date DROP DEFAULT,
    ALTER COLUMN checkin_date TYPE date USING checkin_date::date,
    ALTER COLUMN checkin_date SET DEFAULT CURRENT_DATE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.team_checkin
    ALTER COLUMN checkin_date DROP DEFAULT,
    ALTER COLUMN checkin_date TYPE timestamp with time zone USING (checkin_date::timestamp AT TIME ZONE 'UTC'),
    ALTER COLUMN checkin_date SET DEFAULT CURRENT_DATE;
-- +goose StatementEnd