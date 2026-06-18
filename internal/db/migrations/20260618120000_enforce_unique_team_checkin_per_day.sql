-- +goose Up
-- +goose StatementBegin
WITH ranked_checkins AS (
    SELECT
        id,
        ROW_NUMBER() OVER (
            PARTITION BY team_id, user_id, checkin_date
            ORDER BY created_date DESC, id DESC
        ) AS row_num
    FROM thunderdome.team_checkin
)
DELETE FROM thunderdome.team_checkin tc
USING ranked_checkins rc
WHERE tc.id = rc.id
  AND rc.row_num > 1;

CREATE UNIQUE INDEX IF NOT EXISTS team_checkin_team_id_user_id_checkin_date_uidx
    ON thunderdome.team_checkin USING btree (team_id, user_id, checkin_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS thunderdome.team_checkin_team_id_user_id_checkin_date_uidx;
-- +goose StatementEnd