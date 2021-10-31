--
-- Extensions
--
CREATE extension IF NOT EXISTS "uuid-ossp";

--
-- Tables
--
CREATE TABLE IF NOT EXISTS "alert" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "type" varchar(128) DEFAULT 'NEW',
    "content" text,
    "active" bool DEFAULT true,
    "allow_dismiss" bool DEFAULT true,
    "registered_only" bool DEFAULT true,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "api_keys" (
    "id" text NOT NULL,
    "user_id" uuid NOT NULL,
    "name" varchar(256) NOT NULL,
    "active" bool DEFAULT true,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "battles" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "owner_id" uuid,
    "name" varchar(256),
    "voting_locked" bool DEFAULT true,
    "active_plan_id" uuid,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "point_values_allowed" jsonb DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::jsonb,
    "auto_finish_voting" bool DEFAULT true,
    "point_average_rounding" varchar(5) DEFAULT 'ceil',
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "battles_leaders" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    PRIMARY KEY ("battle_id","user_id")
);

CREATE TABLE IF NOT EXISTS "battles_users" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "active" bool DEFAULT false,
    "abandoned" bool DEFAULT false,
    "spectator" bool DEFAULT false,
    PRIMARY KEY ("battle_id","user_id")
);

CREATE TABLE IF NOT EXISTS "department_team" (
    "department_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","team_id")
);

CREATE TABLE IF NOT EXISTS "department_user" (
    "department_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","user_id")
);

CREATE TABLE IF NOT EXISTS "organization" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "organization_department" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "organization_id" uuid,
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "organization_team" (
    "organization_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","team_id")
);

CREATE TABLE IF NOT EXISTS "organization_user" (
    "organization_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","user_id")
);

CREATE TABLE IF NOT EXISTS "plans" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "points" varchar(3) DEFAULT '',
    "active" bool DEFAULT false,
    "battle_id" uuid NOT NULL,
    "votes" jsonb DEFAULT '[]'::jsonb,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "skipped" bool DEFAULT false,
    "votestart_time" timestamp DEFAULT now(),
    "voteend_time" timestamp DEFAULT now(),
    "acceptance_criteria" text,
    "link" text,
    "description" text,
    "reference_id" varchar(128),
    "type" varchar(64) DEFAULT 'story',
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "team" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "team_battle" (
    "team_id" uuid NOT NULL,
    "battle_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("team_id","battle_id")
);

CREATE TABLE IF NOT EXISTS "team_user" (
    "team_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    PRIMARY KEY ("team_id","user_id")
);

CREATE TABLE IF NOT EXISTS "user_reset" (
    "reset_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '01:00:00'::interval),
    PRIMARY KEY ("reset_id")
);

CREATE TABLE IF NOT EXISTS "user_verify" (
    "verify_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '24:00:00'::interval),
    PRIMARY KEY ("verify_id")
);

CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(64),
    "created_date" timestamp DEFAULT now(),
    "last_active" timestamp DEFAULT now(),
    "email" varchar(320),
    "password" text,
    "type" varchar(128) DEFAULT 'PRIVATE',
    "verified" bool DEFAULT false,
    "avatar" varchar(128) DEFAULT 'identicon',
    "notifications_enabled" bool DEFAULT true,
    "jira_rest_api_token" varchar(128) DEFAULT '',
    "country" varchar(2),
    "company" varchar(256),
    "job_title" varchar(128),
    "updated_date" timestamp DEFAULT now(),
    "locale" varchar(2),
    PRIMARY KEY ("id")
);

--
-- FOREIGN KEYS
--

ALTER TABLE "api_keys" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "battles" ADD FOREIGN KEY ("owner_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "battles_leaders" ADD FOREIGN KEY ("battle_id") REFERENCES "battles"("id") ON DELETE CASCADE;
ALTER TABLE "battles_leaders" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "battles_users" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "battles_users" ADD FOREIGN KEY ("battle_id") REFERENCES "battles"("id") ON DELETE CASCADE;
ALTER TABLE "department_team" ADD FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE;
ALTER TABLE "department_team" ADD FOREIGN KEY ("department_id") REFERENCES "organization_department"("id") ON DELETE CASCADE;
ALTER TABLE "department_user" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "department_user" ADD FOREIGN KEY ("department_id") REFERENCES "organization_department"("id") ON DELETE CASCADE;
ALTER TABLE "organization_department" ADD FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE;
ALTER TABLE "organization_team" ADD FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE;
ALTER TABLE "organization_team" ADD FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE;
ALTER TABLE "organization_user" ADD FOREIGN KEY ("organization_id") REFERENCES "organization"("id") ON DELETE CASCADE;
ALTER TABLE "organization_user" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "plans" ADD FOREIGN KEY ("battle_id") REFERENCES "battles"("id") ON DELETE CASCADE;
ALTER TABLE "team_battle" ADD FOREIGN KEY ("battle_id") REFERENCES "battles"("id") ON DELETE CASCADE;
ALTER TABLE "team_battle" ADD FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE;
ALTER TABLE "team_user" ADD FOREIGN KEY ("team_id") REFERENCES "team"("id") ON DELETE CASCADE;
ALTER TABLE "team_user" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "user_reset" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;
ALTER TABLE "user_verify" ADD FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE CASCADE;

--
-- Types (used in Stored Procedures)
--
DROP TYPE IF EXISTS UsersVote;
CREATE TYPE UsersVote AS
(
    "warriorId"     uuid,
    "vote"   VARCHAR(3)
);

--
-- Views
--
CREATE MATERIALIZED VIEW IF NOT EXISTS active_countries AS SELECT DISTINCT country FROM users;

--
-- Stored Procedures
--

-- Reset All Users to Inactive, used by server restart --
CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_users SET active = false WHERE active = true;
END;
$$;

-- Create a Battle Plan --
CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (id, battle_id, name, type, reference_id, link, description, acceptance_criteria)
    VALUES (planId, battleId, planName, planType, referenceId, planLink, planDescription, acceptanceCriteria);

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Revise Plan --
CREATE OR REPLACE PROCEDURE revise_plan(planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
DECLARE battleId UUID;
BEGIN
    UPDATE plans
    SET
        updated_date = NOW(),
        name = planName,
        type = planType,
        reference_id = referenceId,
        link = planLink,
        description = planDescription,
        acceptance_criteria = acceptanceCriteria
    WHERE id = planId RETURNING battle_id INTO battleId;

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Activate a Battles Plan, and de-activate any current active plan
CREATE OR REPLACE PROCEDURE activate_plan_voting(battleId UUID, planId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false WHERE battle_id = battle_id;
    -- set PlanID active to true
    UPDATE plans SET updated_date = NOW(), active = true, skipped = false, points = '', votestart_time = NOW(), votes = '[]'::jsonb WHERE id = planId;
    -- set battle VotingLocked and ActivePlanID
    UPDATE battles SET updated_date = NOW(), voting_locked = false, active_plan_id = planId WHERE id = battleId;
    COMMIT;
END;
$$;

-- Skip a Battles Plan Voting --
CREATE OR REPLACE PROCEDURE skip_plan_voting(battleId UUID, planId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false, skipped = true, voteend_time = NOW() WHERE battle_id = battleId;
    -- set battle VotingLocked and activePlanId to null
    UPDATE battles SET updated_date = NOW(), voting_locked = true, active_plan_id = null WHERE id = battleId;
    COMMIT;
END;
$$;

-- End a Battles Plan Voting --
CREATE OR REPLACE PROCEDURE end_plan_voting(battleId UUID, planId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false, voteend_time = NOW() WHERE battle_id = battleId;
    -- set battle VotingLocked
    UPDATE battles SET updated_date = NOW(), voting_locked = true WHERE id = battleId;
    COMMIT;
END;
$$;

-- Finalize a plan --
CREATE OR REPLACE PROCEDURE finalize_plan(battleId UUID, planId UUID, planPoints VARCHAR(3))
LANGUAGE plpgsql AS $$
BEGIN
    -- set plan points and deactivate
    UPDATE plans SET updated_date = NOW(), active = false, points = planPoints WHERE id = planId;
    -- reset battle active_plan_id
    UPDATE battles SET updated_date = NOW(), active_plan_id = null WHERE id = battleId;
    COMMIT;
END;
$$;

-- Delete a plan --
CREATE OR REPLACE PROCEDURE delete_plan(battleId UUID, planId UUID)
LANGUAGE plpgsql AS $$
DECLARE active_plan_id UUID;
BEGIN
    active_plan_id := (SELECT b.active_plan_id FROM battles b WHERE b.id = battleId);
    DELETE FROM plans WHERE id = planId;
    
    IF active_plan_id = planId THEN
        UPDATE battles SET updated_date = NOW(), voting_locked = true, active_plan_id = null WHERE id = battleId;
    ELSE
        UPDATE battles SET updated_date = NOW() WHERE id = battleId;
    END IF;
    
    COMMIT;
END;
$$;

-- Set Battle Leader --
CREATE OR REPLACE PROCEDURE set_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
END;
$$;

-- Demote Battle Leader --
CREATE OR REPLACE PROCEDURE demote_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles_leaders WHERE battle_id = battleId AND user_id = leaderId;
END;
$$;

-- Delete Battle --
CREATE OR REPLACE PROCEDURE delete_battle(battleId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles WHERE id = battleId;

    COMMIT;
END;
$$;

-- Set User Vote --
CREATE OR REPLACE PROCEDURE set_user_vote(planId UUID, userId UUID, userVote VARCHAR(3))
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT json_agg(data)
        FROM (
            SELECT coalesce(newVote."warriorId", oldVote."warriorId") AS "warriorId", coalesce(newVote.vote, oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::UsersVote,p1.votes) AS oldVote
            FULL JOIN jsonb_populate_recordset(null::UsersVote,
                ('[{"warriorId":"'|| userId::TEXT ||'", "vote":"'|| userVote ||'"}]')::JSONB
            ) AS newVote
            ON newVote."warriorId" = oldVote."warriorId"
        ) data
    )
    WHERE p1.id = planId;
    
    COMMIT;
END;
$$;

-- Retract User Vote --
CREATE OR REPLACE PROCEDURE retract_user_vote(planId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT coalesce(json_agg(data), '[]'::JSON)
        FROM (
            SELECT coalesce(oldVote."warriorId") AS "warriorId", coalesce(oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::UsersVote,p1.votes) AS oldVote
            WHERE oldVote."warriorId" != userId
        ) data
    )
    WHERE p1.id = planId;
    
    COMMIT;
END;
$$;

-- Reset User Password --
CREATE OR REPLACE PROCEDURE reset_user_password(resetId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM user_reset wr
        LEFT JOIN users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase reset record expired
        DELETE FROM user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;

-- Update User Password --
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Verify a user account email
CREATE OR REPLACE PROCEDURE verify_user_account(verifyId UUID)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM user_verify wv
        LEFT JOIN users w ON w.id = wv.user_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase verify record expired
        DELETE FROM user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$$;

-- Promote User to GENERAL type (ADMIN) by ID --
CREATE OR REPLACE PROCEDURE promote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'GENERAL', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Promote User to GENERAL type (ADMIN) by Email --
CREATE OR REPLACE PROCEDURE promote_user_by_email(userEmail VARCHAR(320))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'GENERAL', updated_date = NOW() WHERE email = userEmail;

    COMMIT;
END;
$$;

-- Demote User to CORPORAL type (Registered) by ID --
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'CORPORAL', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Clean up Battles older than X Days --
CREATE OR REPLACE PROCEDURE clean_battles(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;

-- Clean up Guest Users (and their created battles) older than X Days --
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE last_active < (NOW() - daysOld * interval '1 day') AND type = 'PRIVATE';
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Deletes a user and all his battle(s), api keys --
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Updates a users profile --
CREATE OR REPLACE PROCEDURE user_profile_update(
    userId UUID,
    userName VARCHAR(64),
    userAvatar VARCHAR(128),
    notificationsEnabled BOOLEAN,
    userCountry VARCHAR(2),
    userLocale VARCHAR(2),
    userCompany VARCHAR(256),
    userJobTitle VARCHAR(128)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users
    SET
        name = userName,
        avatar = userAvatar,
        notifications_enabled = notificationsEnabled,
        country = userCountry,
        locale = userLocale,
        company = userCompany,
        job_title = userJobTitle,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;
END;
$$;

--
-- Stored Functions
--

-- Find and update all user emails that include an uppercase character and dont have a duplicate account to lowercase email --
CREATE OR REPLACE FUNCTION lowercase_unique_user_emails() RETURNS table (
    name VARCHAR(256), email VARCHAR(320)
) AS $$
BEGIN
    RETURN QUERY
        UPDATE users u
        SET email = lower(u.email), updated_date = NOW()
        FROM (
            SELECT lower(su.email) AS email
            FROM users su
            WHERE su.email IS NOT NULL
            GROUP BY lower(su.email) HAVING count(su.*) = 1
        ) AS subquery
        WHERE lower(u.email) = subquery.email AND u.email ~ '[A-Z]' RETURNING u.name, u.email;
END;
$$ LANGUAGE plpgsql;

-- Find and merge duplicate email user accounts caused by case senstitive bug --
CREATE OR REPLACE FUNCTION merge_nonunique_user_accounts() RETURNS table (
    name VARCHAR(256), email VARCHAR(320)
) AS $$
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
        FROM users su
        WHERE su.email IS NOT NULL
        GROUP BY lower(su.email) HAVING count(su.*) > 1
        ORDER BY active_date DESC
    LOOP
        -- update battles
        UPDATE battles SET owner_id = usr.id[1] WHERE owner_id = usr.id[2];
        -- update battle_users
        BEGIN
            UPDATE battles_users SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in battle';
        END;
        -- update battle_leaders
        BEGIN
            UPDATE battles_leaders SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update organization_user
        BEGIN
            UPDATE organization_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update department_user
        BEGIN
            UPDATE department_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in department';
        END;
        -- update team_user
        BEGIN
            UPDATE team_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in team';
        END;
        -- delete extra user
        DELETE FROM users u WHERE u.id = usr.id[2];
        -- update merged user
        UPDATE users u SET
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

    -- update active_countries
    REFRESH MATERIALIZED VIEW active_countries;
END;
$$ LANGUAGE plpgsql;

-- Get Application Stats e.g. total user and battle counts
CREATE OR REPLACE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER,
    OUT apikey_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO battle_count FROM battles;
    SELECT COUNT(*) INTO plan_count FROM plans;
    SELECT COUNT(*) INTO organization_count FROM organization;
    SELECT COUNT(*) INTO department_count FROM organization_department;
    SELECT COUNT(*) INTO team_count FROM team;
    SELECT COUNT(*) INTO apikey_count FROM api_keys;
END;
$$ LANGUAGE plpgsql;

-- Insert a new user password reset
CREATE OR REPLACE FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
)
AS $$
BEGIN
    SELECT id, name INTO userId, userName FROM users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Register a new user
CREATE OR REPLACE FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    INSERT INTO users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Register a new user from existing private
CREATE OR REPLACE FUNCTION register_existing_user(
    IN activeUserId UUID,
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    UPDATE users
    SET
        name = userName,
        email = userEmail,
        password = hashedPassword,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Create Battle --
CREATE OR REPLACE FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    OUT battleId UUID
) AS $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding) VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding) RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
END;
$$ LANGUAGE plpgsql;

-- Add Battle Leaders by Emails --
CREATE OR REPLACE FUNCTION add_battle_leaders_by_email(
    IN battleId UUID,
    IN leaderEmails TEXT,
    OUT leaders JSONB
) AS $$
DECLARE
    emails TEXT[];
    leaderEmail TEXT;
BEGIN
    select into emails regexp_split_to_array(leaderEmails,',');
    FOREACH leaderEmail IN ARRAY emails
    LOOP
        INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, (
            SELECT id FROM users WHERE email = leaderEmail
        ));
    END LOOP;

    SELECT CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END
    FROM battles_leaders bl WHERE bl.battle_id = battleId INTO leaders;
END;
$$ LANGUAGE plpgsql;

-- Get a list of countries
CREATE OR REPLACE FUNCTION countries_active() RETURNS table (
    country VARCHAR(2)
) AS $$
BEGIN
    RETURN QUERY SELECT ac.country FROM active_countries ac;
END;
$$ LANGUAGE plpgsql;

-- Get API Keys --
CREATE OR REPLACE FUNCTION apikeys_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id text, name VARCHAR(256), email VARCHAR(320), active BOOLEAN, created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT apk.id, apk.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM api_keys apk
		LEFT JOIN users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Registered Users list --
CREATE OR REPLACE FUNCTION registered_users_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), avatar VARCHAR(128), verified BOOLEAN, country VARCHAR(2), company VARCHAR(256), job_title VARCHAR(128)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, '')
		FROM users u
		WHERE u.email IS NOT NULL
		ORDER BY u.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

--
-- ORGANIZATIONS --
--

-- Get Organization --
CREATE OR REPLACE FUNCTION organization_get_by_id(
    IN orgId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        WHERE o.id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization User Role --
CREATE OR REPLACE FUNCTION organization_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    OUT role VARCHAR(16)
) AS $$
BEGIN
    SELECT ou.role INTO role
    FROM organization_user ou
    WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations --
CREATE OR REPLACE FUNCTION organization_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        ORDER BY o.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations by User --
CREATE OR REPLACE FUNCTION organization_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP, role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM organization_user ou
        LEFT JOIN organization o ON ou.organization_id = o.id
        WHERE ou.user_id = userId
        ORDER BY created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization --
CREATE OR REPLACE FUNCTION organization_create(
    IN userId UUID,
    IN orgName VARCHAR(256),
    OUT organizationId UUID
) AS $$
BEGIN
    INSERT INTO organization (name) VALUES (orgName) RETURNING id INTO organizationId;
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
END;
$$ LANGUAGE plpgsql;

-- Add User to Organization --
CREATE OR REPLACE FUNCTION organization_user_add(
    IN orgId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (orgId, userId, userRole);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Organization --
CREATE OR REPLACE PROCEDURE organization_user_remove(orgId UUID, userId UUID)
AS $$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM organization_department WHERE organization_id = orgId
    LOOP
        CALL department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT ot.team_id
        FROM organization_team ot
        WHERE ot.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM organization_user WHERE organization_id = orgId AND user_id = userId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Users --
CREATE OR REPLACE FUNCTION organization_user_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), ou.role
        FROM organization_user ou
        LEFT JOIN users u ON ou.user_id = u.id
        WHERE ou.organization_id = orgId
        ORDER BY ou.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Teams --
CREATE OR REPLACE FUNCTION organization_team_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM organization_team ot
        LEFT JOIN team t ON ot.team_id = t.id
        WHERE ot.organization_id = orgId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Team --
CREATE OR REPLACE FUNCTION organization_team_create(
    IN orgId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO organization_team (organization_id, team_id) VALUES (orgId, teamId);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Team User Role --
CREATE OR REPLACE FUNCTION organization_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

--
-- DEPARTMENTS --
--

-- Get Department --
CREATE OR REPLACE FUNCTION department_get_by_id(
    IN departmentId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT od.id, od.name, od.created_date, od.updated_date
        FROM organization_department od
        WHERE od.id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department User Role --
CREATE OR REPLACE FUNCTION department_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Departments --
CREATE OR REPLACE FUNCTION department_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT d.id, d.name, d.created_date, d.updated_date
        FROM organization_department d
        WHERE d.organization_id = orgId
        ORDER BY d.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Department --
CREATE OR REPLACE FUNCTION department_create(
    IN orgId UUID,
    IN departmentName VARCHAR(256),
    OUT departmentId UUID
) AS $$
BEGIN
    INSERT INTO organization_department (name, organization_id) VALUES (departmentName, orgId) RETURNING id INTO departmentId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Teams --
CREATE OR REPLACE FUNCTION department_team_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM department_team dt
        LEFT JOIN team t ON dt.team_id = t.id
        WHERE dt.department_id = departmentId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Department Team --
CREATE OR REPLACE FUNCTION department_team_create(
    IN departmentId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO department_team (department_id, team_id) VALUES (departmentId, teamId);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Team User Role --
CREATE OR REPLACE FUNCTION department_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Users --
CREATE OR REPLACE FUNCTION department_user_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), du.role
        FROM department_user du
        LEFT JOIN users u ON du.user_id = u.id
        WHERE du.department_id = departmentId
        ORDER BY du.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add User to Department --
CREATE OR REPLACE FUNCTION department_user_add(
    IN departmentId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
DECLARE orgId UUID;
BEGIN    
    SELECT organization_id INTO orgId FROM organization_user WHERE user_id = userId;

    IF orgId IS NULL THEN
        RAISE EXCEPTION 'User not in Organization -> %', userId USING HINT = 'Please add user to Organization before department';
    END IF;

    INSERT INTO department_user (department_id, user_id, role) VALUES (departmentId, userId, userRole);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Department --
CREATE OR REPLACE PROCEDURE department_user_remove(departmentId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT dt.team_id
        FROM department_team dt
        WHERE dt.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM department_user WHERE department_id = departmentId AND user_id = userId;
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

--
-- TEAMS --
--

-- Get Team --
CREATE OR REPLACE FUNCTION team_get_by_id(
    IN teamId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM team o
        WHERE o.id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Team User Role --
CREATE OR REPLACE FUNCTION team_get_user_role(
    IN userId UUID,
    IN teamId UUID
) RETURNS table (
    role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT tu.role
        FROM team_user tu
        WHERE tu.team_id = teamId AND tu.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Teams --
CREATE OR REPLACE FUNCTION team_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team t
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Teams by User --
CREATE OR REPLACE FUNCTION team_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team_user tu
        LEFT JOIN team t ON tu.team_id = t.id
        WHERE tu.user_id = userId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Team --
CREATE OR REPLACE FUNCTION team_create(
    IN userId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, 'ADMIN');
END;
$$ LANGUAGE plpgsql;

-- Get Team Users --
CREATE OR REPLACE FUNCTION team_user_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), tu.role
        FROM team_user tu
        LEFT JOIN users u ON tu.user_id = u.id
        WHERE tu.team_id = teamId
        ORDER BY tu.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add User to Team --
CREATE OR REPLACE FUNCTION team_user_add(
    IN teamId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, userRole);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Team --
CREATE OR REPLACE PROCEDURE team_user_remove(teamId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user WHERE team_id = teamId AND user_id = userId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Team Battles --
CREATE OR REPLACE FUNCTION team_battle_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name
        FROM team_battle tb
        LEFT JOIN battles b ON tb.battle_id = b.id
        WHERE tb.team_id = teamId
        ORDER BY tb.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add Battle to Team --
CREATE OR REPLACE FUNCTION team_battle_add(
    IN teamId UUID,
    IN battleId UUID
) RETURNS void AS $$
BEGIN
    INSERT INTO team_battle (team_id, battle_id) VALUES (teamId, battleId);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove Battle from Team --
CREATE OR REPLACE FUNCTION team_battle_remove(
    IN teamId UUID,
    IN battleId UUID
) RETURNS void AS $$
BEGIN
    DELETE FROM team_battle WHERE battle_id = battleId AND team_id = teamId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Delete Team --
CREATE OR REPLACE PROCEDURE team_delete(teamId UUID)
AS $$
BEGIN
    DELETE FROM team WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;