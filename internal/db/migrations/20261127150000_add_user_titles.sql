-- +goose Up
-- Descriptive info about a member of an organization, department or team, unrelated to their access control role.
ALTER TABLE thunderdome.organization_user
    ADD COLUMN IF NOT EXISTS title VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notes VARCHAR(1000) NOT NULL DEFAULT '';

ALTER TABLE thunderdome.department_user
    ADD COLUMN IF NOT EXISTS title VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notes VARCHAR(1000) NOT NULL DEFAULT '';

ALTER TABLE thunderdome.team_user
    ADD COLUMN IF NOT EXISTS title VARCHAR(128) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS notes VARCHAR(1000) NOT NULL DEFAULT '';

-- Titles an organization offers for its members. Members store their own copy so custom titles remain possible.
CREATE TABLE IF NOT EXISTS thunderdome.organization_user_title (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS organization_user_title_name_uidx
    ON thunderdome.organization_user_title (organization_id, lower(name));

-- +goose Down
DROP TABLE IF EXISTS thunderdome.organization_user_title;

ALTER TABLE thunderdome.organization_user
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS notes;

ALTER TABLE thunderdome.department_user
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS notes;

ALTER TABLE thunderdome.team_user
    DROP COLUMN IF EXISTS title,
    DROP COLUMN IF EXISTS notes;
