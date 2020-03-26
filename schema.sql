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

--
-- Table Alterations
--
ALTER TABLE battles ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS created_date TIMESTAMP DEFAULT NOW();
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS last_active TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS updated_date TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS skipped BOOL DEFAULT false;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS votestart_time TIMESTAMP DEFAULT NOW();
ALTER TABLE plans ADD COLUMN IF NOT EXISTS voteend_time TIMESTAMP DEFAULT NOW();
ALTER TABLE battles ADD COLUMN IF NOT EXISTS point_values_allowed JSONB DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::JSONB;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS email VARCHAR(320) UNIQUE;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS password TEXT;
ALTER TABLE warriors ADD COLUMN IF NOT EXISTS rank VARCHAR(128) DEFAULT 'PRIVATE';
ALTER TABLE battles ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE plans ALTER COLUMN id SET DEFAULT uuid_generate_v4();
ALTER TABLE warriors ALTER COLUMN id SET DEFAULT uuid_generate_v4();

--
-- Types (used in Stored Procedures)
--
DROP TYPE IF EXISTS WarriorsVote;
CREATE TYPE WarriorsVote AS
(
    "warriorId"     uuid,
    "vote"   VARCHAR(3)
);

--
-- Stored Procedures
--

-- Reset All Warriors to Inactive, used by server restart --
CREATE OR REPLACE PROCEDURE deactivate_all_warriors()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_warriors SET active = false WHERE active = true;
END;
$$;

-- Create a Battle Plan --
CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256))
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (id, battle_id, name) VALUES (planId, battleId, planName);
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

-- Revise Plan Name --
CREATE OR REPLACE PROCEDURE revise_plan_name(planId UUID, planName VARCHAR(256))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE plans SET updated_date = NOW(), name = planName WHERE id = planId;
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
    END IF;
    
    COMMIT;
END;
$$;

-- Set Battle Leader --
CREATE OR REPLACE PROCEDURE set_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles SET updated_date = NOW(), leader_id = leaderId WHERE id = battleId;
END;
$$;

-- Delete Battle --
CREATE OR REPLACE PROCEDURE delete_battle(battleId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM plans WHERE battle_id = battleId;
    DELETE FROM battles_warriors WHERE battle_id = battleId;
    DELETE FROM battles WHERE id = battleId;

    COMMIT;
END;
$$;

-- Set Warrior Vote --
CREATE OR REPLACE PROCEDURE set_warrior_vote(planId UUID, warriorsId UUID, warriorVote VARCHAR(3))
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT json_agg(data)
        FROM (
            SELECT coalesce(newVote."warriorId", oldVote."warriorId") AS "warriorId", coalesce(newVote.vote, oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::WarriorsVote,p1.votes) AS oldVote
            FULL JOIN jsonb_populate_recordset(null::WarriorsVote,
                ('[{"warriorId":"'|| warriorsId::TEXT ||'", "vote":"'|| warriorVote ||'"}]')::JSONB
            ) AS newVote
            ON newVote."warriorId" = oldVote."warriorId"
        ) data
    )
    WHERE p1.id = planId;
    
    COMMIT;
END;
$$;

-- Retract Warrior Vote --
CREATE OR REPLACE PROCEDURE retract_warrior_vote(planId UUID, warriorsId UUID)
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT coalesce(json_agg(data), '[]'::JSON)
        FROM (
            SELECT coalesce(oldVote."warriorId") AS "warriorId", coalesce(oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::WarriorsVote,p1.votes) AS oldVote
            WHERE oldVote."warriorId" != warriorsId
        ) data
    )
    WHERE p1.id = planId;
    
    COMMIT;
END;
$$;

-- Reset Warrior Password --
CREATE OR REPLACE PROCEDURE reset_warrior_password(resetId UUID, warriorPassword TEXT)
LANGUAGE plpgsql AS $$
DECLARE matchedWarriorId UUID;
BEGIN
	matchedWarriorId := (
        SELECT w.id
        FROM warrior_reset wr
        LEFT JOIN warriors w ON w.id = wr.warrior_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedWarriorId IS NULL THEN
        -- attempt delete incase reset record expired
        DELETE FROM warrior_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE warriors SET password = warriorPassword WHERE id = matchedWarriorId;
    DELETE FROM warrior_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;
