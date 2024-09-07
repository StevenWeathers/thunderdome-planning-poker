-- +goose Up
-- +goose StatementBegin

-- Create a new table for retro templates
CREATE TABLE thunderdome.retro_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    format JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE ,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_ownership CHECK (
        (is_public = true AND organization_id IS NULL AND team_id IS NULL) OR
        (is_public = false AND (organization_id IS NOT NULL OR team_id IS NOT NULL))
    )
);

-- Add indexes for better query performance
CREATE INDEX idx_retro_template_organization ON thunderdome.retro_template(organization_id);
CREATE INDEX idx_retro_template_team ON thunderdome.retro_template(team_id);
CREATE INDEX idx_retro_template_public ON thunderdome.retro_template(is_public);

-- Insert pre-defined templates
INSERT INTO thunderdome.retro_template (id, name, code, description, format, is_public)
VALUES
(
    '5c3b4783-82cb-45a4-ac7b-c956c6b4047e', -- hardcoded UUID for the default template
    'Worked/Improvements/Questions',
    'worked_improve_question',
    'Reflect on what worked, what needs improvement, and what questions you have',
    '{"columns": [
                {"name": "worked", "label": "What worked well", "color": "green", "icon": "smiley"},
                {"name": "improve", "label": "What needs improvement", "color": "red", "icon": "frown"},
                {"name": "question", "label": "I want to ask", "color": "blue", "icon": "question"}
        ]}',
    true
),
(
    gen_random_uuid(),
    'Start/Stop/Continue',
    'start_stop_continue',
    'Reflect on what to start doing, stop doing, and continue doing',
    '{"columns": [
            {"name": "start", "label": "Start", "color": "green"},
            {"name": "stop", "label": "Stop", "color": "red"},
            {"name": "continue", "label": "Continue", "color": "blue"}
    ]}',
    true
),
(
    gen_random_uuid(),
    'Mad/Sad/Glad',
    'mad_sad_glad',
    'Express feelings about the sprint: what made you mad, sad, or glad',
    '{"columns": [
                {"name": "mad", "label": "Mad", "color": "red", "icon": "angry"},
                {"name": "sad", "label": "Sad", "color": "blue", "icon": "frown"},
                {"name": "glad", "label": "Glad", "color": "green", "icon": "smiley"}
        ]}',
    true
),
(
    gen_random_uuid(),
    'Drop/Add/Keep/Improve',
    'drop_add_keep_improve',
    'Reflect on what to Drop, Add, Keep, and Improve in the next sprint',
    '{"columns": [
            {"name": "drop", "label": "Drop", "color": "red"},
             {"name": "add", "label": "Add", "color": "green"},
            {"name": "keep", "label": "Keep", "color": "blue"},
            {"name": "improve", "label": "Improve", "color": "yellow"}
    ]}',
    true
),
(
    gen_random_uuid(),
    'Liked/Learned/Lacked/Longed for',
    'liked_learned_lacked_longed_for',
    'Reflect on what was Liked, Learned, Lacked, and Longed for',
    '{"columns": [
            {"name": "liked", "label": "Liked", "color": "green"},
             {"name": "learned", "label": "Learned", "color": "blue"},
            {"name": "lacked", "label": "Lacked", "color": "red"},
            {"name": "longedfor", "label": "Longed for", "color": "yellow"}
    ]}',
    true
);

-- Add a column to the retro table to reference the template used
ALTER TABLE thunderdome.retro
ADD COLUMN template_id UUID REFERENCES thunderdome.retro_template(id) DEFAULT '5c3b4783-82cb-45a4-ac7b-c956c6b4047e';

-- First, drop the existing function
DROP FUNCTION IF EXISTS thunderdome.retro_create(uuid, character varying, character varying, text, text, smallint, character varying, integer, boolean, uuid);

-- Then, create the updated function with the new templateID parameter
CREATE OR REPLACE FUNCTION thunderdome.retro_create(
    userid uuid,
    retroname character varying,
    joincode text,
    facilitatorcode text,
    maxvotes smallint,
    brainstormvisibility character varying,
    phasetimelimitmin integer,
    phaseautoadvance boolean,
    teamid uuid,
    templateid uuid
)
RETURNS uuid
LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (
        owner_id, name, join_code, facilitator_code,
        max_votes, brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
        team_id, template_id
    )
    VALUES (
        userId, retroName, joinCode, facilitatorCode,
        maxVotes, brainstormVisibility, phasetimelimitmin, phaseautoadvance, teamid, templateid
    ) RETURNING id INTO retroId;

    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$function$;

ALTER TABLE thunderdome.retro DROP COLUMN format;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE thunderdome.retro ADD COLUMN format VARCHAR(32) DEFAULT 'worked_improve_question' NOT NULL;

DROP FUNCTION IF EXISTS thunderdome.retro_create(uuid, character varying, text, text, smallint, character varying, integer, boolean, uuid, uuid);

CREATE OR REPLACE FUNCTION thunderdome.retro_create(
    userid uuid,
    retroname character varying,
    fmt character varying,
    joincode text,
    facilitatorcode text,
    maxvotes smallint,
    brainstormvisibility character varying,
    phasetimelimitmin integer,
    phaseautoadvance boolean,
    teamid uuid
)
RETURNS uuid
LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (
        owner_id, name, format, join_code, facilitator_code,
        max_votes, brainstorm_visibility, phase_time_limit_min, phase_auto_advance,
        team_id
    )
    VALUES (
        userId, retroName, fmt, joinCode, facilitatorCode,
        maxVotes, brainstormVisibility, phasetimelimitmin, phaseautoadvance, teamid
    ) RETURNING id INTO retroId;

    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$function$;

-- Remove the template_id column from the retro table
ALTER TABLE thunderdome.retro
DROP COLUMN IF EXISTS template_id;

-- Drop the retro_template table
DROP TABLE IF EXISTS thunderdome.retro_template;

-- +goose StatementEnd
