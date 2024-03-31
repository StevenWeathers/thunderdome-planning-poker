-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.subscription ADD COLUMN team_id uuid REFERENCES thunderdome.team(id) ON DELETE SET NULL;
ALTER TABLE thunderdome.subscription ADD COLUMN organization_id uuid REFERENCES thunderdome.organization(id) ON DELETE SET NULL;
CREATE INDEX IF NOT EXISTS subscription_active_user_id_idx ON thunderdome.subscription USING btree (user_id, active);
CREATE INDEX IF NOT EXISTS subscription_active_team_id_idx ON thunderdome.subscription USING btree (team_id, active);
CREATE INDEX IF NOT EXISTS subscription_active_organization_id_idx ON thunderdome.subscription USING btree (organization_id, active);
ALTER TABLE thunderdome.users DROP COLUMN subscribed;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.subscription DROP COLUMN team_id;
ALTER TABLE thunderdome.subscription DROP COLUMN organization_id;
DROP INDEX IF EXISTS subscription_active_user_id_idx;
DROP INDEX IF EXISTS subscription_active_team_id_idx;
DROP INDEX IF EXISTS subscription_active_organization_id_idx;
ALTER TABLE thunderdome.users ADD COLUMN subscribed boolean NOT NULL DEFAULT false;
-- +goose StatementEnd
