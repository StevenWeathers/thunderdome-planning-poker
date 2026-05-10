-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS thunderdome.color_legend_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    color_legend JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT color_legend_template_check_scope CHECK (
        (is_public = true AND num_nonnulls(organization_id, team_id) = 0) OR
        (is_public = false AND num_nonnulls(organization_id, team_id) = 1)
    )
);

CREATE INDEX IF NOT EXISTS idx_color_legend_template_organization_id
    ON thunderdome.color_legend_template(organization_id);

CREATE INDEX IF NOT EXISTS idx_color_legend_template_team_id
    ON thunderdome.color_legend_template(team_id);

CREATE INDEX IF NOT EXISTS idx_color_legend_template_is_public
    ON thunderdome.color_legend_template(is_public);

DROP FUNCTION IF EXISTS thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text, teamid uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.color_legend_template;

CREATE OR REPLACE FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text, teamid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE storyId UUID;
BEGIN
    INSERT INTO thunderdome.storyboard (owner_id, name, join_code, facilitator_code, team_id)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode, teamid) RETURNING id INTO storyId;
    INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);

    RETURN storyId;
END;
$function$;

-- +goose StatementEnd