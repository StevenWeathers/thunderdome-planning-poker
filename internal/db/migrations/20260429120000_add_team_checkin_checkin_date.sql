-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.team_checkin
    ADD COLUMN IF NOT EXISTS checkin_date timestamp with time zone;

UPDATE thunderdome.team_checkin
SET checkin_date = created_date
WHERE checkin_date IS NULL;

ALTER TABLE thunderdome.team_checkin
    ALTER COLUMN checkin_date SET DEFAULT CURRENT_DATE,
    ALTER COLUMN checkin_date SET NOT NULL;

CREATE INDEX IF NOT EXISTS team_checkin_team_id_checkin_date_idx
    ON thunderdome.team_checkin USING btree (team_id, checkin_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS thunderdome.team_checkin_team_id_checkin_date_idx;

ALTER TABLE thunderdome.team_checkin
    DROP COLUMN IF EXISTS checkin_date;
-- +goose StatementEnd