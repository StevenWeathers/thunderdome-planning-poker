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

CREATE TABLE IF NOT EXISTS warriors (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS plans (
    id UUID NOT NULL PRIMARY KEY,
    name VARCHAR(256),
    points VARCHAR(3) DEFAULT '',
    active BOOL DEFAULT false,
    battle_id UUID REFERENCES battles(id) NOT NULL,
    votes JSONB DEFAULT '[]'::JSONB
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

--
-- Table Alterations
--
ALTER TABLE battles ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS point_values_allowed JSONB DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::JSONB;
ALTER TABLE battles ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS auto_finish_voting BOOL DEFAULT true;
ALTER TABLE battles ADD COLUMN IF NOT EXISTS point_average_rounding VARCHAR(5) DEFAULT 'ceil';

ALTER TABLE warriors ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS last_active TIMESTAMP DEFAULT NOW();
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS email VARCHAR(320) UNIQUE;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS password TEXT;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS rank VARCHAR(128) DEFAULT 'PRIVATE';
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS verified BOOL DEFAULT false;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS avatar VARCHAR(128) DEFAULT 'identicon';
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS notifications_enabled BOOL DEFAULT true;
ALTER TABLE warriors ALTER COLUMN id SET DEFAULT uuid_generate_v4();

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

ALTER TABLE battles_warriors ADD COLUMN IF NOT EXISTS abandoned BOOL DEFAULT false;

DO $$
BEGIN
    -- migrate battles leaderId into new association table
    DECLARE battleLeadersExists TEXT := (SELECT to_regclass('battles_leaders'));
    BEGIN
        IF battleLeadersExists IS NULL THEN
            CREATE TABLE battles_leaders AS SELECT id AS battle_id, leader_id AS warrior_id FROM battles;
            ALTER TABLE battles_leaders ADD PRIMARY KEY (battle_id, warrior_id);
            ALTER TABLE battles_leaders ADD CONSTRAINT bl_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE;
            ALTER TABLE battles_leaders ADD CONSTRAINT bl_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
            ALTER TABLE battles DROP CONSTRAINT IF EXISTS b_leader_id_fkey;
            ALTER TABLE battles RENAME COLUMN leader_id TO owner_id;
            ALTER TABLE battles ADD CONSTRAINT b_owner_id_fkey FOREIGN KEY (owner_id) REFERENCES warriors(id) ON DELETE CASCADE;
        END IF;
    END;

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
        WHEN duplicate_object THEN RAISE NOTICE 'battles_warriors constraint bw_battle_id_fkey already exists';
    END;

    BEGIN
        ALTER TABLE battles_warriors ADD CONSTRAINT bw_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
        EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE 'battles_warriors constraint bw_warrior_id_fkey already exists';
    END;

    BEGIN
        ALTER TABLE plans ADD CONSTRAINT p_battle_id_fkey FOREIGN KEY (battle_id) REFERENCES battles(id) ON DELETE CASCADE;
        EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE 'plans constraint p_battle_id_fkey already exists';
    END;

    BEGIN
        ALTER TABLE warrior_reset ADD CONSTRAINT wr_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
        EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE 'warrior_reset constraint wr_warrior_id_fkey already exists';
    END;

    BEGIN
        ALTER TABLE warrior_verify ADD CONSTRAINT wv_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
        EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE 'warrior_verify constraint wv_warrior_id_fkey already exists';
    END;

    BEGIN
        ALTER TABLE api_keys ADD CONSTRAINT apk_warrior_id_fkey FOREIGN KEY (warrior_id) REFERENCES warriors(id) ON DELETE CASCADE;
        EXCEPTION
        WHEN duplicate_object THEN RAISE NOTICE 'api_keys constraint apk_warrior_id_fkey already exists';
    END;
END $$;

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

-- Reset All Warriors to Inactive, used by server restart --
DROP PROCEDURE IF EXISTS deactivate_all_warriors();
CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_warriors SET active = false WHERE active = true;
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
    INSERT INTO battles_leaders (battle_id, warrior_id) VALUES (battleId, leaderId);
END;
$$;

-- Demote Battle Leader --
CREATE OR REPLACE PROCEDURE demote_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles_leaders WHERE battle_id = battleId AND warrior_id = leaderId;
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

-- Set Warrior Vote --
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

-- Retract Warrior Vote --
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

-- Reset Warrior Password --
DROP PROCEDURE IF EXISTS reset_warrior_password(resetId UUID, warriorPassword TEXT);
CREATE OR REPLACE PROCEDURE reset_user_password(resetId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM warrior_reset wr
        LEFT JOIN warriors w ON w.id = wr.warrior_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase reset record expired
        DELETE FROM warrior_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE warriors SET password = userPassword, last_active = NOW() WHERE id = matchedUserId;
    DELETE FROM warrior_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;

-- Update Warrior Password --
DROP PROCEDURE IF EXISTS update_warrior_password(warriorId UUID, warriorPassword TEXT);
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE warriors SET password = userPassword, last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Verify a warrior account email
DROP PROCEDURE IF EXISTS verify_warrior_account(verifyId UUID);
CREATE OR REPLACE PROCEDURE verify_user_account(verifyId UUID)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM warrior_verify wv
        LEFT JOIN warriors w ON w.id = wv.warrior_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase verify record expired
        DELETE FROM warrior_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE warriors SET verified = 'TRUE', last_active = NOW() WHERE id = matchedUserId;
    DELETE FROM warrior_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$$;

-- Promote Warrior to GENERAL Rank (ADMIN) by ID --
DROP PROCEDURE IF EXISTS promote_warrior(warriorId UUID);
CREATE OR REPLACE PROCEDURE promote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE warriors SET rank = 'GENERAL' WHERE id = userId;

    COMMIT;
END;
$$;

-- Promote Warrior to GENERAL Rank (ADMIN) by Email --
DROP PROCEDURE IF EXISTS promote_warrior_by_email(warriorEmail VARCHAR(320));
CREATE OR REPLACE PROCEDURE promote_user_by_email(userEmail VARCHAR(320))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE warriors SET rank = 'GENERAL' WHERE email = userEmail;

    COMMIT;
END;
$$;

-- Demote Warrior to CORPORAL Rank (Registered) by ID --
DROP PROCEDURE IF EXISTS demote_user(warriorId UUID);
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE warriors SET rank = 'CORPORAL' WHERE id = userId;

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

-- Clean up Guest Warriors (and their created battles) older than X Days --
DROP PROCEDURE IF EXISTS clean_guest_warriors(daysOld INTEGER);
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM warriors WHERE last_active < (NOW() - daysOld * interval '1 day') AND rank = 'PRIVATE';

    COMMIT;
END;
$$;

-- Deletes a warrior and all his battle(s), api keys --
DROP PROCEDURE IF EXISTS delete_warrior(warriorId UUID);
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM warriors WHERE id = userId;

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
CREATE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM warriors WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM warriors WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO battle_count FROM battles;
    SELECT COUNT(*) INTO plan_count FROM plans;
END;
$$ LANGUAGE plpgsql;

-- Insert a new warrior password reset
DROP FUNCTION IF EXISTS insert_warrior_reset(VARCHAR);
DROP FUNCTION IF EXISTS insert_warrior_reset(
    IN warriorEmail VARCHAR(320),
    OUT resetId UUID,
    OUT warriorId UUID,
    OUT warriorName VARCHAR(64)
);
CREATE FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
)
AS $$
BEGIN
    SELECT id, name INTO userId, userName FROM warriors WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO warrior_reset (warrior_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Register a new warrior
DROP FUNCTION IF EXISTS register_warrior(VARCHAR, VARCHAR, TEXT, VARCHAR);
DROP FUNCTION IF EXISTS register_warrior(
    IN warriorName VARCHAR(64),
    IN warriorEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN warriorRank VARCHAR(128),
    OUT warriorId UUID,
    OUT verifyId UUID
);
CREATE FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    INSERT INTO warriors (name, email, password, rank)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO warrior_verify (warrior_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Register a new warrior from existing private
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
CREATE FUNCTION register_existing_user(
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
    UPDATE warriors
    SET
         name = userName,
         email = userEmail,
         password = hashedPassword,
         rank = userType,
         last_active = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO warrior_verify (warrior_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Create Battle --
DROP FUNCTION IF EXISTS create_battle(UUID, VARCHAR, JSONB, BOOL);
DROP FUNCTION IF EXISTS create_battle(UUID, VARCHAR, JSONB, BOOL, VARCHAR);
CREATE FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    OUT battleId UUID
)
AS $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding) VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding) RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, warrior_id) VALUES (battleId, leaderId);
END;
$$ LANGUAGE plpgsql;