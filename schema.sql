--
-- Extensions
--
CREATE extension IF NOT EXISTS "uuid-ossp";

--
-- Tables
--
CREATE TABLE IF NOT EXISTS battles (
    id UUID NOT NULL PRIMARY KEY,
    leader_id UUID,
    name VARCHAR(256),
    voting_locked BOOL DEFAULT true,
    active_plan_id UUID
);

CREATE TABLE IF NOT EXISTS plans (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(256),
    points VARCHAR(3) DEFAULT '',
    active BOOL DEFAULT false,
    battle_id UUID REFERENCES battles(id) NOT NULL,
    votes JSONB DEFAULT '[]'::JSONB
);

--
-- Table Alterations
--
ALTER TABLE battles ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS point_values_allowed JSONB DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::JSONB;
ALTER TABLE battles ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS auto_finish_voting BOOL DEFAULT true;
ALTER TABLE battles ADD COLUMN IF NOT EXISTS point_average_rounding VARCHAR(5) DEFAULT 'ceil';

ALTER TABLE plans ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS skipped BOOL DEFAULT false;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS votestart_time TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS voteend_time TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS link TEXT;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS acceptance_criteria TEXT;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS reference_id VARCHAR(128);
ALTER TABLE plans ADD COLUMN IF NOT EXISTS type VARCHAR(64) DEFAULT 'story';

DO $$
BEGIN
    -- migrate battles leaderId into new association table
    DECLARE battleLeadersExists TEXT := (SELECT to_regclass('battles_leaders'));
    DECLARE usersExists TEXT := (SELECT to_regclass('users'));
    BEGIN
        IF usersExists IS NULL THEN
            CREATE TABLE IF NOT EXISTS warriors (
                id UUID NOT NULL PRIMARY KEY,
                name VARCHAR(64)
            );

            CREATE TABLE IF NOT EXISTS battles_warriors (
                battle_id UUID REFERENCES battles NOT NULL,
                warrior_id UUID REFERENCES warriors NOT NULL,
                active BOOL DEFAULT false,
                PRIMARY KEY (battle_id, warrior_id)
            );

            CREATE TABLE IF NOT EXISTS warrior_reset (
                reset_id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                warrior_id UUID REFERENCES warriors NOT NULL,
                created_date TIMESTAMP DEFAULT NOW(),
                expire_date TIMESTAMP DEFAULT NOW() + INTERVAL '1 hour'
            );

            CREATE TABLE IF NOT EXISTS warrior_verify (
                verify_id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
                warrior_id UUID REFERENCES warriors NOT NULL,
                created_date TIMESTAMP DEFAULT NOW(),
                expire_date TIMESTAMP DEFAULT NOW() + INTERVAL '24 hour'
            );

            CREATE TABLE IF NOT EXISTS api_keys (
                id TEXT NOT NULL PRIMARY KEY,
                warrior_id UUID REFERENCES warriors NOT NULL,
                name VARCHAR(256) NOT NULL,
                active BOOL DEFAULT true,
                created_date TIMESTAMP DEFAULT NOW(),
                updated_date TIMESTAMP DEFAULT NOW(),
                UNIQUE(warrior_id, name)
            );

            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS last_active TIMESTAMP DEFAULT NOW();
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS email VARCHAR(320) UNIQUE;
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS password TEXT;
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS rank VARCHAR(128) DEFAULT 'PRIVATE';
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS verified BOOL DEFAULT false;
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS avatar VARCHAR(128) DEFAULT 'identicon';
            ALTER TABLE warriors ADD COLUMN IF NOT EXISTS notifications_enabled BOOL DEFAULT true;
            ALTER TABLE warriors ALTER COLUMN id SET DEFAULT uuid_generate_v4();

            ALTER TABLE battles_warriors ADD COLUMN IF NOT EXISTS abandoned BOOL DEFAULT false;

            --
            -- Constraints
            --
            ALTER TABLE battles_warriors DROP CONSTRAINT IF EXISTS battles_warriors_battle_id_fkey;
            ALTER TABLE battles_warriors DROP CONSTRAINT IF EXISTS battles_warriors_warrior_id_fkey;
            ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_warrior_id_fkey;
            ALTER TABLE plans DROP CONSTRAINT IF EXISTS plans_battle_id_fkey;
            ALTER TABLE warrior_verify DROP CONSTRAINT IF EXISTS warrior_verify_warrior_id_fkey;
            ALTER TABLE warrior_reset DROP CONSTRAINT IF EXISTS warrior_reset_warrior_id_fkey;

            BEGIN
                ALTER TABLE battles_warriors ADD CONSTRAINT bw_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            BEGIN
                ALTER TABLE battles_warriors ADD CONSTRAINT bw_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            BEGIN
                ALTER TABLE plans ADD CONSTRAINT p_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            BEGIN
                ALTER TABLE warrior_reset ADD CONSTRAINT wr_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            BEGIN
                ALTER TABLE warrior_verify ADD CONSTRAINT wv_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            BEGIN
                ALTER TABLE api_keys ADD CONSTRAINT apk_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
                EXCEPTION
                WHEN duplicate_object THEN
            END;

            ALTER TABLE battles_warriors RENAME TO battles_users;
            ALTER TABLE battles_users RENAME COLUMN warrior_id TO user_id;
            ALTER TABLE warriors RENAME TO users;
            ALTER TABLE users RENAME rank TO type;
            ALTER TABLE warrior_reset RENAME TO user_reset;
            ALTER TABLE user_reset RENAME COLUMN warrior_id TO user_id;
            ALTER TABLE warrior_verify RENAME TO user_verify;
            ALTER TABLE user_verify RENAME COLUMN warrior_id TO user_id;
            ALTER TABLE api_keys RENAME COLUMN warrior_id to user_id;
        END IF;

        IF battleLeadersExists IS NULL THEN
            CREATE TABLE battles_leaders AS SELECT id AS battle_id, leader_id AS warrior_id FROM battles;
            ALTER TABLE battles_leaders ADD PRIMARY KEY (battle_id, warrior_id);
            ALTER TABLE battles_leaders ADD CONSTRAINT bl_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE;
            ALTER TABLE battles_leaders ADD CONSTRAINT bl_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES users(id) ON DELETE CASCADE;
            ALTER TABLE battles DROP CONSTRAINT IF EXISTS b_leader_id_fkey;
            ALTER TABLE battles RENAME COLUMN leader_id TO owner_id;
            ALTER TABLE battles ADD CONSTRAINT b_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;
        END IF;

        
        BEGIN
            ALTER TABLE battles_leaders RENAME COLUMN warrior_id TO user_id;
            EXCEPTION
                WHEN undefined_column THEN
        END;
    END;
END $$;

CREATE TABLE IF NOT EXISTS organization (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS organization_user (
    organization_id UUID,
    user_id UUID,
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (organization_id, user_id),
    CONSTRAINT ou_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT ou_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization_department (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    organization_id UUID,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    UNIQUE(organization_id, name),
    CONSTRAINT od_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS department_user (
    department_id UUID,
    user_id UUID,
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (department_id, user_id),
    CONSTRAINT du_department_id FOREIGN KEY(department_id) REFERENCES organization_department(id) ON DELETE CASCADE,
    CONSTRAINT du_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS team (
    id UUID NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(256),
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_user (
    team_id UUID,
    user_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    role VARCHAR(16) NOT NULL DEFAULT 'MEMBER',
    PRIMARY KEY (team_id, user_id),
    CONSTRAINT tu_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE,
    CONSTRAINT tu_user_id FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS organization_team (
    organization_id UUID,
    team_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (organization_id, team_id),
    UNIQUE(team_id),
    CONSTRAINT ot_organization_id FOREIGN KEY(organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT ot_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS department_team (
    department_id UUID,
    team_id UUID,
    created_date TIMESTAMP DEFAULT NOW(),
    updated_date TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (department_id, team_id),
    UNIQUE(team_id),
    CONSTRAINT dt_department_id FOREIGN KEY(department_id) REFERENCES organization_department(id) ON DELETE CASCADE,
    CONSTRAINT dt_team_id FOREIGN KEY(team_id) REFERENCES team(id) ON DELETE CASCADE
);

--
-- Types (used in Stored Procedures)
--
DROP TYPE IF EXISTS WarriorsVote;
DROP TYPE IF EXISTS UsersVote;
CREATE TYPE UsersVote AS
(
    "warriorId"     uuid,
    "vote"   VARCHAR(3)
);

--
-- Stored Procedures
--

-- Reset All Users to Inactive, used by server restart --
DROP PROCEDURE IF EXISTS deactivate_all_warriors();
CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_users SET active = false WHERE active = true;
END;
$$;

-- Create a Battle Plan --
DROP PROCEDURE IF EXISTS create_plan(battleid uuid, planid uuid, planname character varying);
DROP PROCEDURE IF EXISTS create_plan(battleid uuid, planid uuid, planname character varying, referenceid character varying, planlink text, plandescription text, acceptancecriteria text);
CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (id, battle_id, name, type, reference_id, link, description, acceptance_criteria)
    VALUES (planId, battleId, planName, planType, referenceId, planLink, planDescription, acceptanceCriteria);

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Revise Plan --
DROP PROCEDURE IF EXISTS revise_plan(planid uuid, planname character varying, referenceid character varying, planlink text, plandescription text, acceptancecriteria text);
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

-- Revise Plan Name (Replaced by revise_plan) --
DROP PROCEDURE IF EXISTS revise_plan_name(planId UUID, planName VARCHAR(256));

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
DROP PROCEDURE IF EXISTS set_warrior_vote(planId UUID, warriorsId UUID, warriorVote VARCHAR(3));
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
DROP PROCEDURE IF EXISTS retract_warrior_vote(planId UUID, warriorsId UUID);
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
DROP PROCEDURE IF EXISTS reset_warrior_password(resetId UUID, warriorPassword TEXT);
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

    UPDATE users SET password = userPassword, last_active = NOW() WHERE id = matchedUserId;
    DELETE FROM user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;

-- Update User Password --
DROP PROCEDURE IF EXISTS update_warrior_password(warriorId UUID, warriorPassword TEXT);
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET password = userPassword, last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Verify a user account email
DROP PROCEDURE IF EXISTS verify_warrior_account(verifyId UUID);
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

    UPDATE users SET verified = 'TRUE', last_active = NOW() WHERE id = matchedUserId;
    DELETE FROM user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$$;

-- Promote User to GENERAL type (ADMIN) by ID --
DROP PROCEDURE IF EXISTS promote_warrior(warriorId UUID);
CREATE OR REPLACE PROCEDURE promote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'GENERAL' WHERE id = userId;

    COMMIT;
END;
$$;

-- Promote User to GENERAL type (ADMIN) by Email --
DROP PROCEDURE IF EXISTS promote_warrior_by_email(warriorEmail VARCHAR(320));
CREATE OR REPLACE PROCEDURE promote_user_by_email(userEmail VARCHAR(320))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'GENERAL' WHERE email = userEmail;

    COMMIT;
END;
$$;

-- Demote User to CORPORAL type (Registered) by ID --
DROP PROCEDURE IF EXISTS demote_warrior(warriorId UUID);
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'CORPORAL' WHERE id = userId;

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
DROP PROCEDURE IF EXISTS clean_guest_warriors(daysOld INTEGER);
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE last_active < (NOW() - daysOld * interval '1 day') AND type = 'PRIVATE';

    COMMIT;
END;
$$;

-- Deletes a user and all his battle(s), api keys --
DROP PROCEDURE IF EXISTS delete_warrior(warriorId UUID);
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE id = userId;

    COMMIT;
END;
$$;

--
-- Stored Functions
--

-- Get Application Stats e.g. total user and battle counts
DROP FUNCTION IF EXISTS get_app_stats();
DROP FUNCTION IF EXISTS get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER
);
CREATE OR REPLACE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO battle_count FROM battles;
    SELECT COUNT(*) INTO plan_count FROM plans;
    SELECT COUNT(*) INTO organization_count FROM organization;
    SELECT COUNT(*) INTO department_count FROM organization_department;
    SELECT COUNT(*) INTO team_count FROM team;
END;
$$ LANGUAGE plpgsql;

-- Insert a new user password reset
DROP FUNCTION IF EXISTS insert_warrior_reset(VARCHAR);
DROP FUNCTION IF EXISTS insert_warrior_reset(
    IN warriorEmail VARCHAR(320),
    OUT resetId UUID,
    OUT warriorId UUID,
    OUT warriorName VARCHAR(64)
);
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
DROP FUNCTION IF EXISTS register_warrior(VARCHAR, VARCHAR, TEXT, VARCHAR);
DROP FUNCTION IF EXISTS register_warrior(
    IN warriorName VARCHAR(64),
    IN warriorEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN warriorRank VARCHAR(128),
    OUT warriorId UUID,
    OUT verifyId UUID
);
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
DROP FUNCTION IF EXISTS register_existing_warrior(UUID, VARCHAR, VARCHAR, TEXT, VARCHAR);
DROP FUNCTION IF EXISTS register_existing_warrior(
    IN activeWarriorId UUID,
    IN warriorName VARCHAR(64),
    IN warriorEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN warriorRank VARCHAR(128),
    OUT warriorId UUID,
    OUT verifyId UUID
);
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
        last_active = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Create Battle --
DROP FUNCTION IF EXISTS create_battle(UUID, VARCHAR, JSONB, BOOL);
DROP FUNCTION IF EXISTS create_battle(UUID, VARCHAR, JSONB, BOOL, VARCHAR);
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
        SELECT id, name, created_date, updated_date
        FROM organization
        WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization with Role --
CREATE OR REPLACE FUNCTION organization_get_with_role(
    IN userId UUID,
    IN orgId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP, role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM organization_user ou
        LEFT JOIN organization o ON ou.organization_id = o.id
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
        SELECT id, name, created_date, updated_date
        FROM organization
        ORDER BY created_date
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
        SELECT u.id, u.name, u.email, ou.role
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

--
-- DEPARTMENTS --
--

-- Get Department --
CREATE OR REPLACE FUNCTION organization_department_get_by_id(
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

-- Get Department with Role --
CREATE OR REPLACE FUNCTION organization_department_get_with_role(
    IN userId UUID,
    IN departmentId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP, role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM department_user ou
        LEFT JOIN organization_department o ON ou.department_id = o.id
        WHERE ou.department_id = departmentId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Departments --
CREATE OR REPLACE FUNCTION organization_department_list(
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
CREATE OR REPLACE FUNCTION organization_department_create(
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
CREATE OR REPLACE FUNCTION organization_department_team_list(
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
CREATE OR REPLACE FUNCTION organization_department_team_create(
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

-- Get Department Users --
CREATE OR REPLACE FUNCTION organization_department_user_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, u.email, du.role
        FROM department_user du
        LEFT JOIN users u ON du.user_id = u.id
        WHERE du.department_id = departmentId
        ORDER BY du.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add User to Department --
CREATE OR REPLACE FUNCTION organization_department_user_add(
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

--
-- TEAMS --
--

-- Get Teams --
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
        SELECT u.id, u.name, u.email, tu.role
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