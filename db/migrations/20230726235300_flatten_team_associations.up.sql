ALTER TABLE thunderdome.poker ADD COLUMN team_id UUID REFERENCES thunderdome.team (id) ON DELETE CASCADE;
ALTER TABLE thunderdome.retro ADD COLUMN team_id UUID REFERENCES thunderdome.team (id) ON DELETE CASCADE;
ALTER TABLE thunderdome.storyboard ADD COLUMN team_id UUID REFERENCES thunderdome.team (id) ON DELETE CASCADE;
ALTER TABLE thunderdome.team ADD COLUMN organization_id UUID REFERENCES thunderdome.organization (id) ON DELETE CASCADE;
ALTER TABLE thunderdome.team ADD COLUMN department_id UUID REFERENCES thunderdome.organization_department (id) ON DELETE CASCADE;

WITH teampokers AS (
    SELECT team_id, poker_id FROM thunderdome.team_poker
)
UPDATE thunderdome.poker
SET team_id = teampokers.team_id
FROM teampokers
WHERE thunderdome.poker.id = teampokers.poker_id;
CREATE INDEX poker_team_id_idx ON thunderdome.poker(team_id);

WITH teamretros AS (
    SELECT team_id, retro_id FROM thunderdome.team_retro
)
UPDATE thunderdome.retro
SET team_id = teamretros.team_id
FROM teamretros
WHERE thunderdome.retro.id = teamretros.retro_id;
CREATE INDEX retro_team_id_idx ON thunderdome.retro(team_id);

WITH teamstoryboards AS (
    SELECT team_id, storyboard_id FROM thunderdome.team_storyboard
)
UPDATE thunderdome.storyboard
SET team_id = teamstoryboards.team_id
FROM teamstoryboards
WHERE thunderdome.storyboard.id = teamstoryboards.storyboard_id;
CREATE INDEX storyboard_team_id_idx ON thunderdome.storyboard(team_id);

WITH teamdepartments AS (
    SELECT team_id, department_id FROM thunderdome.department_team
)
UPDATE thunderdome.team
SET department_id = teamdepartments.department_id
FROM teamdepartments
WHERE thunderdome.team.id = teamdepartments.team_id;
CREATE INDEX team_department_id_idx ON thunderdome.team(department_id);

WITH teamorganizations AS (
    SELECT team_id, organization_id FROM thunderdome.organization_team
)
UPDATE thunderdome.team
SET organization_id = teamorganizations.organization_id
FROM teamorganizations
WHERE thunderdome.team.id = teamorganizations.team_id;
CREATE INDEX team_organization_id_idx ON thunderdome.team(organization_id);

DROP FUNCTION thunderdome.department_team_create(departmentid uuid, teamname character varying);
DROP FUNCTION thunderdome.organization_team_create(orgid uuid, teamname character varying);
DROP PROCEDURE thunderdome.organization_delete(IN orgid uuid);
DROP PROCEDURE thunderdome.department_delete(IN deptid uuid);
DROP FUNCTION thunderdome.department_create(orgid uuid, departmentname character varying);
DROP FUNCTION thunderdome.team_create_poker(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT pokerid uuid);
DROP FUNCTION thunderdome.team_create_retro(teamid uuid, userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying);
DROP FUNCTION thunderdome.team_create_storyboard(teamid uuid, ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text);
DROP PROCEDURE thunderdome.team_delete(IN teamid uuid);
DROP PROCEDURE thunderdome.poker_user_vote_retract(IN planid uuid, IN userid uuid);
DROP PROCEDURE thunderdome.poker_user_vote_set(IN planid uuid, IN userid uuid, IN uservote character varying);
DROP FUNCTION thunderdome.sb_goals_get(storyboardid uuid);

DROP FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text);
CREATE FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text, teamid uuid)
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

DROP FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying);
CREATE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, teamid uuid)
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

DROP FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT pokerid uuid);
CREATE FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, teamid uuid, OUT pokerid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.poker (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code, team_id)
        VALUES (leaderid, pokername, pointsAllowed, autoVoting, pointAverageRounding, hideVoterIdentity, joinCode, leaderCode, teamid)
        RETURNING id INTO pokerid;
    INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, leaderid);
    INSERT INTO thunderdome.poker_user (poker_id, user_id) VALUES (pokerid, leaderid);
END;
$function$;

CREATE OR REPLACE PROCEDURE thunderdome.department_user_remove(IN departmentid uuid, IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT t.id
        FROM thunderdome.team t
        WHERE t.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.department_user WHERE department_id = departmentId AND user_id = userId;

    COMMIT;
END;
$procedure$;

CREATE OR REPLACE PROCEDURE thunderdome.organization_user_remove(IN orgid uuid, IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM thunderdome.organization_department WHERE organization_id = orgId
    LOOP
        CALL department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT t.id
        FROM thunderdome.team t
        WHERE t.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.organization_user WHERE organization_id = orgId AND user_id = userId;

    COMMIT;
END;
$procedure$;

DROP TABLE thunderdome.team_poker;
DROP TABLE thunderdome.team_retro;
DROP TABLE thunderdome.team_storyboard;
DROP TABLE thunderdome.department_team;
DROP TABLE thunderdome.organization_team;

-- fix storyboard story move proc --
CREATE OR REPLACE PROCEDURE thunderdome.sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE storyboardId UUID;
DECLARE srcGoalId UUID;
DECLARE srcColumnId UUID;
DECLARE srcSortOrder INTEGER;
DECLARE targetSortOrder INTEGER;
BEGIN
    SET CONSTRAINTS thunderdome.storyboard_story_column_id_sort_order_key DEFERRED;
    -- Get Story current details
    SELECT
        storyboard_id, goal_id, column_id, sort_order, name, color, content, created_date
    INTO
        storyboardId, srcGoalId, srcColumnId, srcSortOrder
    FROM thunderdome.storyboard_story WHERE id = storyId;

    -- Get target sort order
    IF placeBefore = '' THEN
        SELECT coalesce(max(sort_order), 0) + 1 INTO targetSortOrder FROM thunderdome.storyboard_story WHERE column_id = columnId;
    ELSE
        SELECT sort_order INTO targetSortOrder FROM thunderdome.storyboard_story WHERE column_id = columnId AND id = placeBefore::UUID;
    END IF;

    -- Remove from source column
    UPDATE thunderdome.storyboard_story SET column_id = columnId, sort_order = 9000 WHERE id = storyId;
    -- Update sort order in src column
    UPDATE thunderdome.storyboard_story ss SET sort_order = (t.sort_order - 1)
    FROM (
        SELECT id, sort_order FROM thunderdome.storyboard_story
        WHERE column_id = srcColumnId AND sort_order > srcSortOrder
        ORDER BY sort_order ASC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Update sort order for any story that should come after newly moved story
    UPDATE thunderdome.storyboard_story ss SET sort_order = (t.sort_order + 1)
    FROM (
        SELECT id, sort_order FROM thunderdome.storyboard_story
        WHERE column_id = columnId AND sort_order >= targetSortOrder
        ORDER BY sort_order DESC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Finally, insert story in its ordered place
	UPDATE thunderdome.storyboard_story SET sort_order = targetSortOrder WHERE id = storyId;

    COMMIT;
END;
$procedure$;
