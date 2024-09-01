-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.retro ADD COLUMN phase_time_limit_min SMALLINT NOT NULL DEFAULT 0;
ALTER TABLE thunderdome.retro ADD COLUMN phase_time_start TIMESTAMPTZ NOT NULL DEFAULT NOW();
DROP FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, teamid uuid);
CREATE OR REPLACE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, phasetimelimitmin int, teamid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (
        owner_id, name, format, join_code, facilitator_code,
        max_votes, brainstorm_visibility, phase_time_limit_min,
        team_id
    )
    VALUES (
        userId, retroName, fmt, joinCode, facilitatorCode,
        maxVotes, brainstormVisibility, phasetimelimitmin, teamid
    ) RETURNING id INTO retroId;
    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$function$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.retro DROP COLUMN phase_time_limit_min;
ALTER TABLE thunderdome.retro DROP COLUMN phase_time_start;
DROP FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, phasetimelimitmin int, teamid uuid);
CREATE OR REPLACE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, teamid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility, team_id)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility, teamid) RETURNING id INTO retroId;
    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$function$;
-- +goose StatementEnd
