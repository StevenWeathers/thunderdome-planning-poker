-- +goose Up
-- +goose StatementBegin

-- Add vote_count column to the existing table
ALTER TABLE "thunderdome"."retro_group_vote"
ADD COLUMN "vote_count" integer NOT NULL DEFAULT 1;

-- Create an index on vote_count for performance
CREATE INDEX "idx_retro_group_vote_vote_count"
ON "thunderdome"."retro_group_vote" ("vote_count");

-- Add allow_cumulative_voting column to the existing table
ALTER TABLE "thunderdome"."retro"
ADD COLUMN "allow_cumulative_voting" boolean NOT NULL DEFAULT false;

DROP FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, phasetimelimitmin integer, phaseautoadvance boolean, teamid uuid, templateid uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop allow_cumulative_voting column from the table
ALTER TABLE "thunderdome"."retro" DROP COLUMN "allow_cumulative_voting";

-- Drop vote_count column from the table
ALTER TABLE "thunderdome"."retro_group_vote"
DROP COLUMN "vote_count";

CREATE OR REPLACE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, phasetimelimitmin integer, phaseautoadvance boolean, teamid uuid, templateid uuid)
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

-- +goose StatementEnd
