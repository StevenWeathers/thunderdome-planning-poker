-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.team_kudos (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    team_id uuid NOT NULL REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    target_user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    comment text,
    kudos_date date DEFAULT CURRENT_DATE NOT NULL,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    updated_date timestamp with time zone DEFAULT now() NOT NULL
);

CREATE INDEX IF NOT EXISTS team_kudos_team_id_idx ON thunderdome.team_kudos USING btree (team_id);
CREATE INDEX IF NOT EXISTS team_kudos_user_id_idx ON thunderdome.team_kudos USING btree (user_id);
CREATE INDEX IF NOT EXISTS team_kudos_target_user_id_idx ON thunderdome.team_kudos USING btree (target_user_id);
CREATE INDEX IF NOT EXISTS team_kudos_kudos_date_idx ON thunderdome.team_kudos USING btree (kudos_date);
CREATE UNIQUE INDEX IF NOT EXISTS team_kudos_daily_unique_idx
    ON thunderdome.team_kudos USING btree (team_id, user_id, target_user_id, kudos_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS thunderdome.team_kudos;
-- +goose StatementEnd