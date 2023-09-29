-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.appstats_get(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT poker_count integer, OUT poker_story_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_poker_count integer, OUT active_poker_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer) RETURNS record
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.department_user_add(departmentid uuid, userid uuid, userrole character varying) RETURNS void
    LANGUAGE plpgsql
    AS $$
DECLARE orgId UUID;
BEGIN
    SELECT organization_id INTO orgId FROM thunderdome.organization_user WHERE user_id = userId;

    IF orgId IS NULL THEN
        RAISE EXCEPTION 'User not in Organization -> %', userId USING HINT = 'Please add user to Organization before department';
    END IF;

    INSERT INTO thunderdome.department_user (department_id, user_id, role) VALUES (departmentId, userId, userRole);
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.department_user_remove(IN departmentid uuid, IN userid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT t.id
        FROM thunderdome.team t
        WHERE t.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.department_user WHERE department_id = departmentId AND user_id = userId;

    COMMIT;
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.organization_create(userid uuid, orgname character varying) RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
    LANGUAGE plpgsql
    AS $$
    DECLARE organizationId uuid;
BEGIN
    INSERT INTO thunderdome.organization (name) VALUES (orgName) RETURNING thunderdome.organization.id INTO organizationId;
    INSERT INTO thunderdome.organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
    RETURN QUERY SELECT o.id, o.name, o.created_date, o.updated_date FROM thunderdome.organization o
        WHERE o.id = organizationID;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.organization_user_remove(IN orgid uuid, IN userid uuid)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, teamid uuid, OUT pokerid uuid) RETURNS uuid
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO thunderdome.poker (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code, team_id)
        VALUES (leaderid, pokername, pointsAllowed, autoVoting, pointAverageRounding, hideVoterIdentity, joinCode, leaderCode, teamid)
        RETURNING id INTO pokerid;
    INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, leaderid);
    INSERT INTO thunderdome.poker_user (poker_id, user_id) VALUES (pokerid, leaderid);
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.poker_facilitator_add_by_email(pokerid uuid, facilitatoremails text, OUT facilitators jsonb) RETURNS jsonb
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE PROCEDURE thunderdome.poker_plan_voting_stop(IN pokerid uuid, IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, voteend_time = NOW()
    WHERE poker_id = pokerid AND id = storyid;
    -- set battle VotingLocked
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), voting_locked = true WHERE id = pokerid;
    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.poker_story_activate(IN pokerid uuid, IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false WHERE poker_id = pokerid AND active = true;
    -- set id active to true
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = true, skipped = false, points = '', votestart_time = NOW(), votes = '[]'::jsonb WHERE id = storyid;
    -- set battle voting_locked and active_story_id
    UPDATE thunderdome.poker SET last_active = NOW(), updated_date = NOW(), voting_locked = false, active_story_id = storyid WHERE id = pokerid;
    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE
	active_storyid UUID;
	story RECORD;
    pos DOUBLE PRECISION = -1;
BEGIN
    active_storyid := (SELECT b.active_story_id FROM thunderdome.poker b WHERE b.id = pokerid);
    DELETE FROM thunderdome.poker_story WHERE id = storyid;

	FOR story IN SELECT id, position FROM thunderdome.poker_story WHERE poker_id = pokerid ORDER BY position
    LOOP
        pos = pos + 1;

        UPDATE thunderdome.poker_story SET position = pos WHERE id = story.id;
    END LOOP;


    IF active_storyid = storyid THEN
        UPDATE thunderdome.poker SET last_active = NOW(), voting_locked = true, active_story_id = null
        WHERE id = pokerid;
    END IF;

    COMMIT;
END$$;

CREATE OR REPLACE PROCEDURE thunderdome.poker_story_finalize(IN pokerid uuid, IN storyid uuid, IN storypoints character varying)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- set points and deactivate
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, points = storypoints WHERE id = storyid;
    -- reset battle active_story_id
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), active_story_id = null WHERE id = pokerid;
    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.poker_vote_skip(IN battleid uuid, IN planid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    -- set current active to false
    UPDATE thunderdome.poker_story SET updated_date = NOW(), active = false, skipped = true, voteend_time = NOW() WHERE poker_id = battleid;
    -- set battle voting_locked and active_story_id to null
    UPDATE thunderdome.poker SET updated_date = NOW(), last_active = NOW(), voting_locked = true, active_story_id = null WHERE id = battleid;
    COMMIT;
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.prune_team_checkins() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.prune_user_sessions() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.refresh_active_countries() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, teamid uuid) RETURNS uuid
    LANGUAGE plpgsql
    AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO thunderdome.retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility, team_id)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility, teamid) RETURNING id INTO retroId;
    INSERT INTO thunderdome.retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO thunderdome.retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_column_delete(IN columnid uuid)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text, teamid uuid) RETURNS uuid
    LANGUAGE plpgsql
    AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO thunderdome.storyboard (owner_id, name, join_code, facilitator_code, team_id)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode, teamid) RETURNING id INTO storyId;
    INSERT INTO thunderdome.storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO thunderdome.storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);

    RETURN storyId;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_goal_delete(IN goalid uuid)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_story_delete(IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.team_create(userid uuid, teamname character varying) RETURNS TABLE(id uuid, name character varying, created_date timestamp with time zone, updated_date timestamp with time zone)
    LANGUAGE plpgsql
    AS $$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO thunderdome.team (name) VALUES (teamName) RETURNING thunderdome.team.id INTO teamId;
    INSERT INTO thunderdome.team_user (team_id, user_id, role) VALUES (teamId, userId, 'ADMIN');
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM thunderdome.team t WHERE t.id = teamId;
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.update_poker_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.update_retro_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.update_storyboard_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE PROCEDURE thunderdome.user_account_verify(IN verifyid uuid)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE PROCEDURE thunderdome.user_disable(IN userid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE thunderdome.users SET disabled = true, updated_date = NOW()
        WHERE id = userId;
    DELETE FROM thunderdome.user_session WHERE user_id = userId;

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_enable(IN userid uuid, IN mfasecret text)
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO thunderdome.user_mfa (user_id, secret) VALUES (userId, mfaSecret);
    UPDATE thunderdome.users SET mfa_enabled = true, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.user_mfa_remove(IN userid uuid)
    LANGUAGE plpgsql
    AS $$
BEGIN
    DELETE FROM thunderdome.user_mfa WHERE user_id = userId;
    UPDATE thunderdome.users SET mfa_enabled = false, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text)
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.user_register(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid) RETURNS record
    LANGUAGE plpgsql
    AS $$
BEGIN
    INSERT INTO thunderdome.users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO thunderdome.user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.user_register_existing(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid) RETURNS record
    LANGUAGE plpgsql
    AS $$
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
$$;

CREATE OR REPLACE FUNCTION thunderdome.user_reset_create(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying) RETURNS record
    LANGUAGE plpgsql
    AS $$
BEGIN
    SELECT id, name INTO userId, userName FROM thunderdome.users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO thunderdome.user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE thunderdome.users_deactivate_all()
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE thunderdome.poker_user SET active = false WHERE active = true;
    UPDATE thunderdome.retro_user SET active = false WHERE active = true;
    UPDATE thunderdome.storyboard_user SET active = false WHERE active = true;
END;
$$;

CREATE OR REPLACE FUNCTION thunderdome.users_registered_email_search(email_search character varying, l_limit integer, l_offset integer) RETURNS TABLE(id uuid, name character varying, email character varying, type character varying, avatar character varying, verified boolean, country character varying, company character varying, job_title character varying, count integer)
    LANGUAGE plpgsql
    AS $$
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
$$;

DROP TRIGGER IF EXISTS poker_facilitator_poker_last_active ON thunderdome.poker_facilitator;
DROP TRIGGER IF EXISTS poker_story_poker_last_active ON thunderdome.poker_story;
DROP TRIGGER IF EXISTS poker_user_poker_last_active ON thunderdome.poker_user;
DROP TRIGGER IF EXISTS prune_team_checkins ON thunderdome.team_checkin;
DROP TRIGGER IF EXISTS prune_user_sessions ON thunderdome.user_session;
DROP TRIGGER IF EXISTS retro_action_retro_last_active ON thunderdome.retro_action;
DROP TRIGGER IF EXISTS retro_facilitator_retro_last_active ON thunderdome.retro_facilitator;
DROP TRIGGER IF EXISTS retro_group_retro_last_active ON thunderdome.retro_group;
DROP TRIGGER IF EXISTS retro_item_retro_last_active ON thunderdome.retro_item;
DROP TRIGGER IF EXISTS retro_user_retro_last_active ON thunderdome.retro_user;
DROP TRIGGER IF EXISTS storyboard_facilitator_storyboard_last_active ON thunderdome.storyboard_facilitator;
DROP TRIGGER IF EXISTS storyboard_story_storyboard_last_active ON thunderdome.storyboard_story;
DROP TRIGGER IF EXISTS storyboard_user_storyboard_last_active ON thunderdome.storyboard_user;
DROP TRIGGER IF EXISTS user_refresh_active_countries ON thunderdome.users;
CREATE TRIGGER poker_facilitator_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_facilitator FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
CREATE TRIGGER poker_story_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_story FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
CREATE TRIGGER poker_user_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_user FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
CREATE TRIGGER prune_team_checkins AFTER INSERT ON thunderdome.team_checkin FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_team_checkins();
CREATE TRIGGER prune_user_sessions AFTER INSERT ON thunderdome.user_session FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_user_sessions();
CREATE TRIGGER retro_action_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_action FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
CREATE TRIGGER retro_facilitator_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_facilitator FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
CREATE TRIGGER retro_group_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_group FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
CREATE TRIGGER retro_item_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_item FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
CREATE TRIGGER retro_user_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_user FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
CREATE TRIGGER storyboard_facilitator_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_facilitator FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
CREATE TRIGGER storyboard_story_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_story FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
CREATE TRIGGER storyboard_user_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_user FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
CREATE TRIGGER user_refresh_active_countries AFTER INSERT OR DELETE OR UPDATE ON thunderdome.users FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.refresh_active_countries();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER poker_facilitator_poker_last_active ON thunderdome.poker_facilitator;
DROP TRIGGER poker_story_poker_last_active ON thunderdome.poker_story;
DROP TRIGGER poker_user_poker_last_active ON thunderdome.poker_user;
DROP TRIGGER prune_team_checkins ON thunderdome.team_checkin;
DROP TRIGGER prune_user_sessions ON thunderdome.user_session;
DROP TRIGGER retro_action_retro_last_active ON thunderdome.retro_action;
DROP TRIGGER retro_facilitator_retro_last_active ON thunderdome.retro_facilitator;
DROP TRIGGER retro_group_retro_last_active ON thunderdome.retro_group;
DROP TRIGGER retro_item_retro_last_active ON thunderdome.retro_item;
DROP TRIGGER retro_user_retro_last_active ON thunderdome.retro_user;
DROP TRIGGER storyboard_facilitator_storyboard_last_active ON thunderdome.storyboard_facilitator;
DROP TRIGGER storyboard_story_storyboard_last_active ON thunderdome.storyboard_story;
DROP TRIGGER storyboard_user_storyboard_last_active ON thunderdome.storyboard_user;
DROP TRIGGER user_refresh_active_countries ON thunderdome.users;

DROP FUNCTION thunderdome.appstats_get(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT poker_count integer, OUT poker_story_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_poker_count integer, OUT active_poker_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer);
DROP FUNCTION thunderdome.department_user_add(departmentid uuid, userid uuid, userrole character varying);
DROP PROCEDURE thunderdome.department_user_remove(IN departmentid uuid, IN userid uuid);
DROP FUNCTION thunderdome.organization_create(userid uuid, orgname character varying);
DROP PROCEDURE thunderdome.organization_user_remove(IN orgid uuid, IN userid uuid);
DROP FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, teamid uuid, OUT pokerid uuid);
DROP FUNCTION thunderdome.poker_facilitator_add_by_email(pokerid uuid, facilitatoremails text, OUT facilitators jsonb);
DROP PROCEDURE thunderdome.poker_plan_voting_stop(IN pokerid uuid, IN storyid uuid);
DROP PROCEDURE thunderdome.poker_story_activate(IN pokerid uuid, IN storyid uuid);
DROP PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid);
DROP PROCEDURE thunderdome.poker_story_finalize(IN pokerid uuid, IN storyid uuid, IN storypoints character varying);
DROP PROCEDURE thunderdome.poker_vote_skip(IN battleid uuid, IN planid uuid);
DROP FUNCTION thunderdome.prune_team_checkins();
DROP FUNCTION thunderdome.prune_user_sessions();
DROP FUNCTION thunderdome.refresh_active_countries();
DROP FUNCTION thunderdome.retro_create(userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying, teamid uuid);
DROP PROCEDURE thunderdome.sb_column_delete(IN columnid uuid);
DROP FUNCTION thunderdome.sb_create(ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text, teamid uuid);
DROP PROCEDURE thunderdome.sb_goal_delete(IN goalid uuid);
DROP PROCEDURE thunderdome.sb_story_delete(IN storyid uuid);
DROP PROCEDURE thunderdome.sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text);
DROP FUNCTION thunderdome.team_create(userid uuid, teamname character varying);
DROP FUNCTION thunderdome.update_poker_last_active();
DROP FUNCTION thunderdome.update_retro_last_active();
DROP FUNCTION thunderdome.update_storyboard_last_active();
DROP PROCEDURE thunderdome.user_account_verify(IN verifyid uuid);
DROP PROCEDURE thunderdome.user_disable(IN userid uuid);
DROP PROCEDURE thunderdome.user_mfa_enable(IN userid uuid, IN mfasecret text);
DROP PROCEDURE thunderdome.user_mfa_remove(IN userid uuid);
DROP PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text);
DROP FUNCTION thunderdome.user_register(username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid);
DROP FUNCTION thunderdome.user_register_existing(activeuserid uuid, username character varying, useremail character varying, hashedpassword text, usertype character varying, OUT userid uuid, OUT verifyid uuid);
DROP FUNCTION thunderdome.user_reset_create(useremail character varying, OUT resetid uuid, OUT userid uuid, OUT username character varying);
DROP PROCEDURE thunderdome.users_deactivate_all();
DROP FUNCTION thunderdome.users_registered_email_search(email_search character varying, l_limit integer, l_offset integer);
-- +goose StatementEnd
