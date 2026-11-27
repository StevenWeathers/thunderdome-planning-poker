-- +goose Up
-- +goose StatementBegin

-- tags are stored as plain text so items keep their labels when a preset is renamed or removed
ALTER TABLE thunderdome.retro_item ADD COLUMN IF NOT EXISTS tags TEXT[] NOT NULL DEFAULT '{}';

-- preset tags offered in retros, either global (public) or owned by an organization or team
CREATE TABLE IF NOT EXISTS thunderdome.retro_tag (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT retro_tag_check_scope CHECK (
        (is_public = true AND num_nonnulls(organization_id, team_id) = 0) OR
        (is_public = false AND num_nonnulls(organization_id, team_id) = 1)
    )
);

-- names are unique (ignoring case) within each scope
CREATE UNIQUE INDEX IF NOT EXISTS idx_retro_tag_public_name
    ON thunderdome.retro_tag(lower(name)) WHERE is_public = true;

CREATE UNIQUE INDEX IF NOT EXISTS idx_retro_tag_organization_name
    ON thunderdome.retro_tag(organization_id, lower(name)) WHERE organization_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_retro_tag_team_name
    ON thunderdome.retro_tag(team_id, lower(name)) WHERE team_id IS NOT NULL;

INSERT INTO thunderdome.retro_tag (name, is_public)
VALUES
    ('Not team responsibility', true),
    ('Needs discussion', true),
    ('Blocked', true),
    ('Quick win', true);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.retro_tag;
ALTER TABLE thunderdome.retro_item DROP COLUMN IF EXISTS tags;

-- +goose StatementEnd
