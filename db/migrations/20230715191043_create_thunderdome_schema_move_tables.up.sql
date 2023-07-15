CREATE SCHEMA thunderdome;

-- rename battle to poker
ALTER TABLE public.battles RENAME TO poker;
ALTER TABLE public.poker RENAME COLUMN active_plan_id TO active_story_id;
ALTER TABLE public.battles_users RENAME TO poker_user;
ALTER TABLE public.poker_user RENAME COLUMN battle_id TO poker_id;
ALTER TABLE public.battles_leaders RENAME TO poker_facilitator;
ALTER TABLE public.poker_facilitator RENAME COLUMN battle_id TO poker_id;
ALTER TABLE public.plans RENAME TO poker_story;
ALTER TABLE public.poker_story RENAME COLUMN battle_id TO poker_id;
ALTER TABLE public.team_battle RENAME TO team_poker;
ALTER TABLE public.team_poker RENAME COLUMN battle_id TO poker_id;
-- rename api_keys to api_key
ALTER TABLE public.api_keys RENAME TO api_key;
-- move tables to thunderdome schema
ALTER TYPE public.UsersVote SET SCHEMA thunderdome;
ALTER MATERIALIZED VIEW public.active_countries SET SCHEMA thunderdome;
ALTER TABLE public.alert SET SCHEMA thunderdome;
ALTER TABLE public.api_key SET SCHEMA thunderdome;
ALTER TABLE public.poker_user SET SCHEMA thunderdome;
ALTER TABLE public.poker_facilitator SET SCHEMA thunderdome;
ALTER TABLE public.poker SET SCHEMA thunderdome;
ALTER TABLE public.department_team SET SCHEMA thunderdome;
ALTER TABLE public.department_user SET SCHEMA thunderdome;
ALTER TABLE public.organization SET SCHEMA thunderdome;
ALTER TABLE public.organization_department SET SCHEMA thunderdome;
ALTER TABLE public.organization_team SET SCHEMA thunderdome;
ALTER TABLE public.organization_user SET SCHEMA thunderdome;
ALTER TABLE public.poker_story SET SCHEMA thunderdome;
ALTER TABLE public.retro SET SCHEMA thunderdome;
ALTER TABLE public.retro_action SET SCHEMA thunderdome;
ALTER TABLE public.retro_action_assignee SET SCHEMA thunderdome;
ALTER TABLE public.retro_action_comment SET SCHEMA thunderdome;
ALTER TABLE public.retro_facilitator SET SCHEMA thunderdome;
ALTER TABLE public.retro_group SET SCHEMA thunderdome;
ALTER TABLE public.retro_group_vote SET SCHEMA thunderdome;
ALTER TABLE public.retro_item SET SCHEMA thunderdome;
ALTER TABLE public.retro_user SET SCHEMA thunderdome;
ALTER TABLE public.storyboard SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_column SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_column_persona SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_facilitator SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_goal SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_goal_persona SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_persona SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_story SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_story_comment SET SCHEMA thunderdome;
ALTER TABLE public.storyboard_user SET SCHEMA thunderdome;
ALTER TABLE public.team SET SCHEMA thunderdome;
ALTER TABLE public.team_poker SET SCHEMA thunderdome;
ALTER TABLE public.team_checkin SET SCHEMA thunderdome;
ALTER TABLE public.team_checkin_comment SET SCHEMA thunderdome;
ALTER TABLE public.team_retro SET SCHEMA thunderdome;
ALTER TABLE public.team_storyboard SET SCHEMA thunderdome;
ALTER TABLE public.team_user SET SCHEMA thunderdome;
ALTER TABLE public.user_mfa SET SCHEMA thunderdome;
ALTER TABLE public.user_reset SET SCHEMA thunderdome;
ALTER TABLE public.user_session SET SCHEMA thunderdome;
ALTER TABLE public.user_verify SET SCHEMA thunderdome;
ALTER TABLE public.users SET SCHEMA thunderdome;
-- update timestamps to with timezone --
ALTER TABLE thunderdome.alert ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.alert ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.api_key ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.api_key ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.department_team ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.department_team ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.department_user ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.department_user ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_department ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_department ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_team ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_team ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_user ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.organization_user ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker_story ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker_story ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker_story ALTER votestart_time TYPE timestamptz USING votestart_time AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.poker_story ALTER voteend_time TYPE timestamptz USING voteend_time AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team_poker ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team_poker ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team_user ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.team_user ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_reset ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_reset ALTER expire_date TYPE timestamptz USING expire_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_session ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_session ALTER expire_date TYPE timestamptz USING expire_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_verify ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.user_verify ALTER expire_date TYPE timestamptz USING expire_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.users ALTER created_date TYPE timestamptz USING created_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.users ALTER updated_date TYPE timestamptz USING updated_date AT TIME ZONE 'UTC';
ALTER TABLE thunderdome.users ALTER last_active TYPE timestamptz USING last_active AT TIME ZONE 'UTC';

-- create triggers
CREATE OR REPLACE FUNCTION thunderdome.refresh_active_countries()
  RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       REFRESH MATERIALIZED VIEW thunderdome.active_countries;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       IF NEW.country <> OLD.country THEN
            REFRESH MATERIALIZED VIEW thunderdome.active_countries;
        END IF;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        REFRESH MATERIALIZED VIEW thunderdome.active_countries;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER user_refresh_active_countries AFTER INSERT OR UPDATE OR DELETE ON thunderdome.users
 EXECUTE PROCEDURE thunderdome.refresh_active_countries();

ALTER TABLE thunderdome.poker ADD COLUMN last_active TIMESTAMPTZ DEFAULT NOW();
UPDATE thunderdome.poker SET last_active = updated_date;
ALTER TABLE thunderdome.poker ALTER COLUMN last_active SET NOT NULL;
CREATE OR REPLACE FUNCTION thunderdome.update_poker_last_active()
  RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.poker SET last_active = NOW() WHERE id = OLD.poker_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.poker SET last_active = NOW() WHERE id = NEW.poker_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.poker SET last_active = NOW() WHERE id = NEW.poker_id;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER poker_user_poker_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.poker_user
 EXECUTE PROCEDURE thunderdome.update_poker_last_active();
CREATE TRIGGER poker_facilitator_poker_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.poker_facilitator
 EXECUTE PROCEDURE thunderdome.update_poker_last_active();
CREATE TRIGGER poker_story_poker_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.poker_story
 EXECUTE PROCEDURE thunderdome.update_poker_last_active();

ALTER TABLE thunderdome.storyboard ADD COLUMN last_active TIMESTAMPTZ DEFAULT NOW();
UPDATE thunderdome.storyboard SET last_active = updated_date;
ALTER TABLE thunderdome.storyboard ALTER COLUMN last_active SET NOT NULL;
CREATE OR REPLACE FUNCTION thunderdome.update_storyboard_last_active()
  RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = OLD.storyboard_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = NEW.storyboard_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = NEW.storyboard_id;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER storyboard_user_storyboard_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.storyboard_user
 EXECUTE PROCEDURE thunderdome.update_storyboard_last_active();
CREATE TRIGGER storyboard_facilitator_storyboard_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.storyboard_facilitator
 EXECUTE PROCEDURE thunderdome.update_storyboard_last_active();
CREATE TRIGGER storyboard_story_storyboard_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.storyboard_story
 EXECUTE PROCEDURE thunderdome.update_storyboard_last_active();

ALTER TABLE thunderdome.retro ADD COLUMN last_active TIMESTAMPTZ DEFAULT NOW();
UPDATE thunderdome.retro SET last_active = updated_date;
ALTER TABLE thunderdome.retro ALTER COLUMN last_active SET NOT NULL;
CREATE OR REPLACE FUNCTION thunderdome.update_retro_last_active()
  RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.retro SET last_active = NOW() WHERE id = OLD.retro_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.retro SET last_active = NOW() WHERE id = NEW.retro_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.retro SET last_active = NOW() WHERE id = NEW.retro_id;
        RETURN NEW;
    END IF;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER retro_user_retro_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.retro_user
 EXECUTE PROCEDURE thunderdome.update_retro_last_active();
CREATE TRIGGER retro_facilitator_retro_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.retro_facilitator
 EXECUTE PROCEDURE thunderdome.update_retro_last_active();
CREATE TRIGGER retro_group_retro_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.retro_group
 EXECUTE PROCEDURE thunderdome.update_retro_last_active();
CREATE TRIGGER retro_item_retro_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.retro_item
 EXECUTE PROCEDURE thunderdome.update_retro_last_active();
CREATE TRIGGER retro_action_retro_last_active AFTER INSERT OR UPDATE OR DELETE ON thunderdome.retro_action
 EXECUTE PROCEDURE thunderdome.update_retro_last_active();
--
-- updated funcs and procs to thunderdome schema
--
DROP PROCEDURE public.activate_plan_voting(IN battleid uuid, IN planid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.poker_story_activate(IN pokerid uuid, IN storyid uuid)
LANGUAGE plpgsql
AS $procedure$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false WHERE poker_id = pokerid AND active = true;
    -- set id active to true
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = true, skipped = false, points = '', votestart_time = NOW(), votes = '[]'::jsonb WHERE id = storyid;
    -- set battle voting_locked and active_story_id
    UPDATE thunderdome.poker SET last_active = NOW(), updated_date = NOW(), voting_locked = false, active_story_id = storyid WHERE id = pokerid;
    COMMIT;
END;
$procedure$;

DROP FUNCTION public.add_battle_leaders_by_email(battleid uuid, leaderemails text, OUT leaders jsonb);
CREATE OR REPLACE FUNCTION thunderdome.poker_facilitator_add_by_email(pokerid uuid, facilitatoremails text, OUT facilitators jsonb)
 RETURNS jsonb
 LANGUAGE plpgsql
AS $function$
DECLARE
    emails TEXT[];
    userEmail TEXT;
BEGIN
    select into emails regexp_split_to_array(facilitatoremails,',');
    FOREACH userEmail IN ARRAY emails
    LOOP
        INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, (
            SELECT id FROM thunderdome.users WHERE LOWER(email) = userEmail
        ));
    END LOOP;

    SELECT CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END
    FROM thunderdome.poker_facilitator bl WHERE bl.poker_id = pokerid INTO facilitators;
END;
$function$;

DROP FUNCTION public.create_battle(leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT battleid uuid);
CREATE OR REPLACE FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT pokerid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.poker (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code)
        VALUES (leaderid, pokername, pointsAllowed, autoVoting, pointAverageRounding, hideVoterIdentity, joinCode, leaderCode)
        RETURNING id INTO pokerid;
    INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, leaderid);
    INSERT INTO thunderdome.poker_user (poker_id, user_id) VALUES (pokerid, leaderid);
END;
$function$;

DROP FUNCTION public.create_retro(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying);
CREATE OR REPLACE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$function$;

DROP FUNCTION public.create_storyboard(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text);
CREATE OR REPLACE FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE storyId UUID;
BEGIN
    INSERT INTO thunderdome.storyboard (owner_id, name, join_code, facilitator_code)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode) RETURNING id INTO storyId;
    INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);

    RETURN storyId;
END;
$function$;

DROP PROCEDURE public.deactivate_all_users();
CREATE OR REPLACE PROCEDURE thunderdome.users_deactivate_all()
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    UPDATE thunderdome.poker_user SET active = false WHERE active = true;
    UPDATE thunderdome.retro_user SET active = false WHERE active = true;
    UPDATE thunderdome.storyboard_user SET active = false WHERE active = true;
END;
$procedure$;

DROP PROCEDURE public.delete_plan(IN battleid uuid, IN planid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE active_storyid UUID;
BEGIN
    active_storyid := (SELECT b.active_story_id FROM thunderdome.poker b WHERE b.id = pokerid);
    DELETE FROM thunderdome.poker_story WHERE id = storyid;

    IF active_storyid = storyid THEN
        UPDATE thunderdome.poker SET last_active = NOW(), voting_locked = true, active_story_id = null
        WHERE id = pokerid;
    END IF;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.delete_storyboard_column(IN columnid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.sb_column_delete(IN columnid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE goalId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT goal_id, sort_order INTO goalId, sortOrder FROM thunderdome.storyboard_column WHERE id = columnId;

    DELETE FROM thunderdome.storyboard_story WHERE column_id = columnId;
    DELETE FROM thunderdome.storyboard_column WHERE id = columnId RETURNING storyboard_id INTO storyboardId;
    UPDATE thunderdome.storyboard_column sc SET sort_order = (sc.sort_order - 1)
        WHERE sc.goal_id = goalId AND sc.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.delete_storyboard_goal(IN goalid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.sb_goal_delete(IN goalid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE storyboardId UUID;
DECLARE sortOrder INTEGER;
BEGIN
    SELECT sort_order, storyboard_id INTO sortOrder, storyboardId FROM thunderdome.storyboard_goal WHERE id = goalId;

    DELETE FROM thunderdome.storyboard_story WHERE goal_id = goalId;
    DELETE FROM thunderdome.storyboard_column WHERE goal_id = goalId;
    DELETE FROM thunderdome.storyboard_goal WHERE id = goalId;
    UPDATE thunderdome.storyboard_goal sg SET sort_order = (sg.sort_order - 1)
        WHERE sg.storyboard_id = storyBoardId AND sg.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.delete_storyboard_story(IN storyid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.sb_story_delete(IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE columnId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT column_id, sort_order, storyboard_id INTO columnId, sortOrder, storyboardId
        FROM thunderdome.storyboard_story WHERE id = storyId;
    DELETE FROM thunderdome.storyboard_story WHERE id = storyId;
    UPDATE thunderdome.storyboard_story ss SET sort_order = (ss.sort_order - 1)
        WHERE ss.column_id = columnId AND ss.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

DROP FUNCTION public.department_create(orgid uuid, departmentname character varying);
CREATE OR REPLACE FUNCTION thunderdome.department_create(orgid uuid, departmentname character varying)
 RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
 LANGUAGE plpgsql
AS $function$
    DECLARE departmentId uuid;
BEGIN
    INSERT INTO thunderdome.organization_department (name, organization_id) VALUES (departmentName, orgId) RETURNING thunderdome.organization_department.id INTO departmentId;
    RETURN QUERY SELECT d.id, d.name, d.created_date, d.updated_date FROM thunderdome.organization_department d
        WHERE d.id = departmentId;
END;
$function$;

DROP PROCEDURE public.department_delete(IN deptid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.department_delete(IN deptid uuid)
 LANGUAGE plpgsql
AS $procedure$
    DECLARE t record;
BEGIN
    FOR t IN SELECT team_id FROM thunderdome.department_team WHERE department_id = deptId
    LOOP
	    CALL thunderdome.team_delete(t.team_id);
    END LOOP;

    DELETE FROM thunderdome.organization_department WHERE id = deptId;

    COMMIT;
END;
$procedure$;

DROP FUNCTION public.department_team_create(departmentid uuid, teamname character varying);
CREATE OR REPLACE FUNCTION thunderdome.department_team_create(departmentid uuid, teamname character varying)
 RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
 LANGUAGE plpgsql
AS $function$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO thunderdome.team (name) VALUES (teamName) RETURNING thunderdome.team.id INTO teamId;
    INSERT INTO thunderdome.department_team (department_id, team_id) VALUES (departmentId, teamId);
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM thunderdome.team t WHERE t.id = teamId;
END;
$function$;

DROP FUNCTION public.department_user_add(departmentid uuid, userid uuid, userrole character varying);
CREATE OR REPLACE FUNCTION thunderdome.department_user_add(departmentid uuid, userid uuid, userrole character varying)
 RETURNS void
 LANGUAGE plpgsql
AS $function$
DECLARE orgId UUID;
BEGIN
    SELECT organization_id INTO orgId FROM thunderdome.organization_user WHERE user_id = userId;

    IF orgId IS NULL THEN
        RAISE EXCEPTION 'User not in Organization -> %', userId USING HINT = 'Please add user to Organization before department';
    END IF;

    INSERT INTO thunderdome.department_user (department_id, user_id, role) VALUES (departmentId, userId, userRole);
END;
$function$;

DROP PROCEDURE public.department_user_remove(IN departmentid uuid, IN userid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.department_user_remove(IN departmentid uuid, IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT dt.team_id
        FROM thunderdome.department_team dt
        WHERE dt.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.department_user WHERE department_id = departmentId AND user_id = userId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.end_plan_voting(IN battleid uuid, IN planid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.poker_plan_voting_stop(IN pokerid uuid, IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, voteend_time = NOW()
    WHERE poker_id = pokerid AND id = storyid;
    -- set battle VotingLocked
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), voting_locked = true WHERE id = pokerid;
    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.finalize_plan(IN battleid uuid, IN planid uuid, IN planpoints character varying);
CREATE OR REPLACE PROCEDURE thunderdome.poker_story_finalize(IN pokerid uuid, IN storyid uuid, IN storypoints character varying)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    -- set points and deactivate
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, points = storypoints WHERE id = storyid;
    -- reset battle active_story_id
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), active_story_id = null WHERE id = pokerid;
    COMMIT;
END;
$procedure$;

DROP FUNCTION public.get_app_stats(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT battle_count integer, OUT plan_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_battle_count integer, OUT active_battle_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer);
CREATE OR REPLACE FUNCTION thunderdome.appstats_get(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT poker_count integer, OUT poker_story_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_poker_count integer, OUT active_poker_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM thunderdome.users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM thunderdome.users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO poker_count FROM thunderdome.poker;
    SELECT COUNT(*) INTO poker_story_count FROM thunderdome.poker_story;
    SELECT COUNT(*) INTO organization_count FROM thunderdome.organization;
    SELECT COUNT(*) INTO department_count FROM thunderdome.organization_department;
    SELECT COUNT(*) INTO team_count FROM thunderdome.team;
    SELECT COUNT(*) INTO apikey_count FROM thunderdome.api_key;
    SELECT COUNT(DISTINCT poker_id), COUNT(user_id)
        INTO active_poker_count, active_poker_user_count
        FROM thunderdome.poker_user WHERE active IS true;
    SELECT COUNT(*) INTO team_checkins_count FROM thunderdome.team_checkin;
    SELECT COUNT(*) INTO retro_count FROM thunderdome.retro;
    SELECT COUNT(DISTINCT retro_id), COUNT(user_id)
        INTO active_retro_count, active_retro_user_count
        FROM thunderdome.retro_user WHERE active IS true;
    SELECT COUNT(*) INTO retro_item_count FROM thunderdome.retro_item;
    SELECT COUNT(*) INTO retro_action_count FROM thunderdome.retro_action;
    SELECT COUNT(*) INTO storyboard_count FROM thunderdome.storyboard;
    SELECT COUNT(DISTINCT storyboard_id), COUNT(user_id)
        INTO active_storyboard_count, active_storyboard_user_count
        FROM thunderdome.storyboard_user WHERE active IS true;
    SELECT COUNT(*) INTO storyboard_goal_count FROM thunderdome.storyboard_goal;
    SELECT COUNT(*) INTO storyboard_column_count FROM thunderdome.storyboard_column;
    SELECT COUNT(*) INTO storyboard_story_count FROM thunderdome.storyboard_story;
    SELECT COUNT(*) INTO storyboard_persona_count FROM thunderdome.storyboard_persona;
END;
$function$;

DROP FUNCTION public.get_storyboard_goals(storyboardid uuid);
CREATE OR REPLACE FUNCTION thunderdome.sb_goals_get(storyboardid uuid)
 RETURNS TABLE(id uuid, sort_order integer, name character varying, columns json, personas json)
 LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
        SELECT
            sg.id,
            sg.sort_order,
            sg.name,
            COALESCE(json_agg(to_jsonb(t) - 'goal_id' ORDER BY t.sort_order) FILTER (WHERE t.id IS NOT NULL), '[]') AS columns,
            COALESCE(json_agg(to_jsonb(sgp) - 'goal_id') FILTER (WHERE sgp.goal_id IS NOT NULL), '[]') AS personas
        FROM thunderdome.storyboard_goal sg
        LEFT JOIN (
            SELECT
                sc.*,
                COALESCE(
                    json_agg(stss ORDER BY stss.sort_order) FILTER (WHERE stss.id IS NOT NULL), '[]'
                ) AS stories,
                COALESCE(
                    json_agg(scp) FILTER (WHERE scp.column_id IS NOT NULL), '[]'
                ) AS personas
            FROM thunderdome.storyboard_column sc
            LEFT JOIN (
                SELECT cp.column_id, sp.*
                FROM thunderdome.storyboard_column_persona cp
                LEFT JOIN thunderdome.storyboard_persona sp ON sp.id = cp.persona_id
            ) scp ON scp.column_id = sc.id
            LEFT JOIN (
                SELECT
                    ss.*,
                    COALESCE(
                        json_agg(stcm ORDER BY stcm.created_date) FILTER (WHERE stcm.id IS NOT NULL), '[]'
                    ) AS comments
                FROM thunderdome.storyboard_story ss
                LEFT JOIN thunderdome.storyboard_story_comment stcm ON stcm.story_id = ss.id
                GROUP BY ss.id
            ) stss ON stss.column_id = sc.id
            GROUP BY sc.id
        ) t ON t.goal_id = sg.id
        LEFT JOIN (
            SELECT gp.goal_id, sp.*
            FROM thunderdome.storyboard_goal_persona gp
            LEFT JOIN thunderdome.storyboard_persona sp ON sp.id = gp.persona_id
        ) sgp ON sgp.goal_id = sg.id
        WHERE sg.storyboard_id = storyboardId
        GROUP BY sg.id
        ORDER BY sg.sort_order;
END;
$function$;

DROP FUNCTION public.insert_user_reset(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying);
CREATE OR REPLACE FUNCTION thunderdome.user_reset_create(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    SELECT id, name INTO userId, userName FROM thunderdome.users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO thunderdome.user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$function$;

DROP FUNCTION public.merge_nonunique_user_accounts();
CREATE OR REPLACE FUNCTION thunderdome.user_merge_nonunique_accounts()
 RETURNS TABLE(name character varying, email character varying)
 LANGUAGE plpgsql
AS $function$
DECLARE usr RECORD;
BEGIN
    FOR usr IN
        SELECT
            array_agg(su.id) as id,
            array_agg(su.name) as name,
            lower(su.email) AS email,
            MAX(last_active) as active_date,
            array_agg(su.country) as country,
            array_agg(su.company) as company,
            array_agg(su.job_title) as job_title,
            array_agg(su.locale) as locale
        FROM thunderdome.users su
        WHERE su.email IS NOT NULL
        GROUP BY lower(su.email) HAVING count(su.*) > 1
        ORDER BY active_date DESC
    LOOP
        -- update poker
        UPDATE thunderdome.poker SET owner_id = usr.id[1] WHERE owner_id = usr.id[2];
        -- update poker_user
        BEGIN
            UPDATE thunderdome.poker_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in poker game';
        END;
        -- update poker_facilitator
        BEGIN
            UPDATE thunderdome.poker_facilitator SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already a poker game facilitator';
        END;
        -- update organization_user
        BEGIN
            UPDATE thunderdome.organization_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update department_user
        BEGIN
            UPDATE thunderdome.department_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in department';
        END;
        -- update team_user
        BEGIN
            UPDATE thunderdome.team_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in team';
        END;
        -- delete extra user
        DELETE FROM thunderdome.users u WHERE u.id = usr.id[2];
        -- update merged user
        UPDATE thunderdome.users u SET
            email = usr.email,
            updated_date = NOW(),
            country = COALESCE(usr.country[1], usr.country[2]),
            company = COALESCE(usr.company[1], usr.company[2]),
            job_title = COALESCE(usr.job_title[1], usr.job_title[2]),
            locale = COALESCE(usr.locale[1], usr.locale[2])
            WHERE u.id = usr.id[1];

        name := usr.name[1];
        email := usr.email;

        RETURN NEXT;
    END LOOP;
END;
$function$;

DROP PROCEDURE public.move_story(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text);
CREATE OR REPLACE PROCEDURE thunderdome.sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE storyboardId UUID;
DECLARE srcGoalId UUID;
DECLARE srcColumnId UUID;
DECLARE srcSortOrder INTEGER;
DECLARE targetSortOrder INTEGER;
BEGIN
    SET CONSTRAINTS thunderdome.storyboard_story.storyboard_story_column_id_sort_order_key DEFERRED;
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

DROP FUNCTION public.organization_create(userid uuid, orgname character varying);
CREATE OR REPLACE FUNCTION thunderdome.organization_create(userid uuid, orgname character varying)
 RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
 LANGUAGE plpgsql
AS $function$
    DECLARE organizationId uuid;
BEGIN
    INSERT INTO thunderdome.organization (name) VALUES (orgName) RETURNING thunderdome.organization.id INTO organizationId;
    INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
    RETURN QUERY SELECT o.id, o.name, o.created_date, o.updated_date FROM thunderdome.organization o
        WHERE o.id = organizationID;
END;
$function$;

DROP PROCEDURE public.organization_delete(IN orgid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.organization_delete(IN orgid uuid)
 LANGUAGE plpgsql
AS $procedure$
    DECLARE d record;
BEGIN
    FOR d IN SELECT id FROM thunderdome.organization_department WHERE organization_id = orgId
    LOOP
	    CALL department_delete(d.id);
    END LOOP;

    DELETE FROM thunderdome.organization WHERE id = orgId;

    COMMIT;
END;
$procedure$;

DROP FUNCTION public.organization_team_create(orgid uuid, teamname character varying);
CREATE OR REPLACE FUNCTION thunderdome.organization_team_create(orgid uuid, teamname character varying)
 RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
 LANGUAGE plpgsql
AS $function$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO thunderdome.team (name) VALUES (teamName) RETURNING thunderdome.team.id INTO teamId;
    INSERT INTO thunderdome.organization_team (organization_id, team_id) VALUES (orgId, teamId);
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM team t WHERE t.id = teamId;
END;
$function$;

DROP PROCEDURE public.organization_user_remove(IN orgid uuid, IN userid uuid);
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
        SELECT ot.team_id
        FROM thunderdome.organization_team ot
        WHERE ot.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.organization_user WHERE organization_id = orgId AND user_id = userId;

    COMMIT;
END;
$procedure$;

ALTER FUNCTION public.prune_team_checkins() SET SCHEMA thunderdome;
CREATE OR REPLACE FUNCTION thunderdome.prune_team_checkins()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.team_checkin WHERE created_date < (NOW() - '60 days'::interval); -- clean up old checkins
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM team_checkin', row_count;
  END IF;
  RETURN NULL;
END;
$function$;

ALTER FUNCTION public.prune_user_sessions() SET SCHEMA thunderdome;
CREATE OR REPLACE FUNCTION thunderdome.prune_user_sessions()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.user_session WHERE expire_date < NOW(); -- clean up old sessions
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM user_session', row_count;
  END IF;
  RETURN NULL;
END;
$function$;

DROP FUNCTION public.register_existing_user(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid);
CREATE OR REPLACE FUNCTION thunderdome.user_register_existing(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    UPDATE thunderdome.users
    SET
        name = userName,
        email = userEmail,
        password = hashedPassword,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$function$;

DROP FUNCTION public.register_user(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid);
CREATE OR REPLACE FUNCTION thunderdome.user_register(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$function$;

DROP FUNCTION public.registered_users_email_search(email_search character varying, l_limit integer, l_offset integer);
CREATE OR REPLACE FUNCTION thunderdome.users_registered_email_search(email_search character varying, l_limit integer, l_offset integer)
 RETURNS TABLE(id uuid, name character varying, email character varying, type character varying, avatar character varying, verified boolean, country character varying, company character varying, job_title character varying, count integer)
 LANGUAGE plpgsql
AS $function$
    DECLARE count INTEGER;
BEGIN
    SELECT count(*)
    FROM thunderdome.users u
    WHERE u.email IS NOT NULL AND u.email ILIKE ('%' || email_search || '%') AND u.disabled IS FALSE INTO count;

    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), count
        FROM thunderdome.users u
        WHERE u.email IS NOT NULL AND u.email ILIKE ('%' || email_search || '%') AND u.disabled IS FALSE
        ORDER BY u.created_date
        LIMIT l_limit
        OFFSET l_offset;
END;
$function$;

DROP PROCEDURE public.reset_user_password(IN resetid uuid, IN userpassword text);
CREATE OR REPLACE PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_reset wr
        LEFT JOIN thunderdome.users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case reset record expired
        DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE thunderdome.users SET password = userPassword, last_active = NOW(), updated_date = NOW()
        WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.retract_user_vote(IN planid uuid, IN userid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.poker_user_vote_retract(IN planid uuid, IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
	UPDATE thunderdome.poker_story p1
    SET votes = (
        SELECT coalesce(json_agg(data), '[]'::JSON)
        FROM (
            SELECT coalesce(oldVote."warriorId") AS "warriorId", coalesce(oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::thunderdome.UsersVote,p1.votes) AS oldVote
            WHERE oldVote."warriorId" != userId
        ) data
    )
    WHERE p1.id = planId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.set_user_vote(IN planid uuid, IN userid uuid, IN uservote character varying);
CREATE OR REPLACE PROCEDURE thunderdome.poker_user_vote_set(IN planid uuid, IN userid uuid, IN uservote character varying)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
	UPDATE thunderdome.poker_story p1
    SET votes = (
        SELECT json_agg(data)
        FROM (
            SELECT coalesce(newVote."warriorId", oldVote."warriorId") AS "warriorId", coalesce(newVote.vote, oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::thunderdome.UsersVote,p1.votes) AS oldVote
            FULL JOIN jsonb_populate_recordset(null::thunderdome.UsersVote,
                ('[{"warriorId":"'|| userId::TEXT ||'", "vote":"'|| userVote ||'"}]')::JSONB
            ) AS newVote
            ON newVote."warriorId" = oldVote."warriorId"
        ) data
    )
    WHERE p1.id = planId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.skip_plan_voting(IN battleid uuid, IN planid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.poker_vote_skip(IN battleid uuid, IN planid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, skipped = true, voteend_time = NOW() WHERE poker_id = battleid;
    -- set battle voting_locked and active_story_id to null
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), voting_locked = true, active_story_id = null WHERE id = battleid;
    COMMIT;
END;
$procedure$;

DROP FUNCTION public.team_create(userid uuid, teamname character varying);
CREATE OR REPLACE FUNCTION thunderdome.team_create(userid uuid, teamname character varying)
 RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
 LANGUAGE plpgsql
AS $function$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO thunderdome.team (name) VALUES (teamName) RETURNING thunderdome.team.id INTO teamId;
    INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES (teamId, userId, 'ADMIN');
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM thunderdome.team t WHERE t.id = teamId;
END;
$function$;

DROP FUNCTION public.team_create_battle(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT battleid uuid);
CREATE OR REPLACE FUNCTION thunderdome.team_create_poker(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, OUT pokerid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.poker (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code)
        VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding, hidevoteridentity, joinCode, leaderCode)
        RETURNING id INTO pokerid;
    INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, leaderId);
    INSERT INTO thunderdome.poker_user (poker_id, user_id) VALUES (pokerid, leaderId);
    INSERT INTO thunderdome.team_poker (team_id, poker_id) VALUES (teamid, pokerid);
END;
$function$;

DROP FUNCTION public.team_create_retro(teamid uuid, userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying);
CREATE OR REPLACE FUNCTION thunderdome.team_create_retro(teamid uuid, userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.team_retro (team_id, retro_id) VALUES (teamid, retroId);

    RETURN retroId;
END;
$function$;

DROP FUNCTION public.team_create_storyboard(teamid uuid, ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text);
CREATE OR REPLACE FUNCTION thunderdome.team_create_storyboard(teamid uuid, ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
DECLARE storyId UUID;
BEGIN
    INSERT INTO thunderdome.storyboard (owner_id, name, join_code, facilitator_code)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode) RETURNING id INTO storyId;
    INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);
    INSERT INTO thunderdome.team_storyboard (team_id, storyboard_id) VALUES (teamid, storyId);

    RETURN storyId;
END;
$function$;

DROP PROCEDURE public.team_delete(IN teamid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.team_delete(IN teamid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.poker WHERE id IN (
        SELECT poker_id FROM team_poker WHERE team_id = teamid
    );

    DELETE FROM thunderdome.retro WHERE id IN (
        SELECT retro_id FROM thunderdome.team_retro WHERE team_id = teamid
    );

    DELETE FROM thunderdome.storyboard WHERE id IN (
        SELECT storyboard_id FROM thunderdome.team_storyboard WHERE team_id = teamid
    );

    DELETE FROM thunderdome.team WHERE id = teamid;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.user_disable(IN userid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.user_disable(IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    UPDATE thunderdome.users SET disabled = true, updated_date = NOW()
        WHERE id = userId;
    DELETE FROM thunderdome.user_session WHERE user_id = userId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.user_mfa_enable(IN userid uuid, IN mfasecret text);
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_enable(IN userid uuid, IN mfasecret text)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    INSERT INTO thunderdome.user_mfa (user_id, secret) VALUES (userId, mfaSecret);
    UPDATE thunderdome.users SET mfa_enabled = true, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.user_mfa_remove(IN userid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_remove(IN userid uuid)
 LANGUAGE plpgsql
AS $procedure$
BEGIN
    DELETE FROM thunderdome.user_mfa WHERE user_id = userId;
    UPDATE thunderdome.users SET mfa_enabled = false, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$procedure$;

DROP PROCEDURE public.verify_user_account(IN verifyid uuid);
CREATE OR REPLACE PROCEDURE thunderdome.user_account_verify(IN verifyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_verify wv
        LEFT JOIN thunderdome.users w ON w.id = wv.user_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case verify record expired
        DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE thunderdome.users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$procedure$;
