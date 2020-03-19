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
    battle_id UUID references battles(id) NOT NULL,
    votes JSONB DEFAULT '[]'::JSONB
);

CREATE TABLE IF NOT EXISTS battles_warriors (
    battle_id UUID references battles NOT NULL,
    warrior_id UUID REFERENCES warriors NOT NULL,
    active BOOL DEFAULT false,
    PRIMARY KEY (battle_id, warrior_id)
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

--
-- Stored Procedures
--

-- Reset All Warriors to Inactive, used by server restart --
DROP PROCEDURE IF EXISTS deactivate_all_warriors();
CREATE PROCEDURE deactivate_all_warriors()
LANGUAGE SQL AS $$
    UPDATE battles_warriors SET active = false WHERE active = true;
$$;

-- Reset All Warriors to Inactive, used by server restart --
DROP PROCEDURE IF EXISTS retreat_warrior(battle_id UUID, warrior_id UUID);
CREATE PROCEDURE retreat_warrior(battle_id UUID, warrior_id UUID)
LANGUAGE SQL AS $$
    -- set warrior to inactive in battle_id
    UPDATE battles_warriors SET active = false WHERE battle_id = battle_id AND warrior_id = warrior_id;
    -- set warrior last_active
    UPDATE warriors SET last_active = NOW() WHERE id = warrior_id;
$$;

-- Activate a Battles Plan, and de-activate any current active plan
DROP PROCEDURE IF EXISTS activate_plan_voting(battle_id UUID, plan_id UUID);
CREATE PROCEDURE activate_plan_voting(battle_id UUID, plan_id UUID)
LANGUAGE SQL AS $$
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false WHERE battle_id = battle_id;
    -- set PlanID active to true
    UPDATE plans SET updated_date = NOW(), active = true, skipped = false, points = '', votestart_time = NOW(), votes = '[]'::jsonb WHERE id = plan_id;
    -- set battle VotingLocked and ActivePlanID
    UPDATE battles SET updated_date = NOW(), voting_locked = false, active_plan_id = plan_id WHERE id = battle_id;
$$;

-- Skip a Battles Plan Voting --
DROP PROCEDURE IF EXISTS skip_plan_voting(battle_id UUID, plan_id UUID);
CREATE PROCEDURE skip_plan_voting(battle_id UUID, plan_id UUID)
LANGUAGE SQL AS $$
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false, skipped = true, voteend_time = NOW() WHERE battle_id = battle_id;
    -- set battle VotingLocked and activePlanId to null
    UPDATE battles SET updated_date = NOW(), voting_locked = true, active_plan_id = null WHERE id = battle_id;
$$;

-- End a Battles Plan Voting --
DROP PROCEDURE IF EXISTS end_plan_voting(battle_id UUID, plan_id UUID);
CREATE PROCEDURE end_plan_voting(battle_id UUID, plan_id UUID)
LANGUAGE SQL AS $$
    -- set current active to false
    UPDATE plans SET updated_date = NOW(), active = false, voteend_time = NOW() WHERE battle_id = battle_id;
    -- set battle VotingLocked
    UPDATE battles SET updated_date = NOW(), voting_locked = true WHERE id = battle_id;
$$;

-- Finalize a plan --
DROP PROCEDURE IF EXISTS finalize_plan(battle_id UUID, plan_id UUID, plan_points VARCHAR(3));
CREATE PROCEDURE finalize_plan(battle_id UUID, plan_id UUID, plan_points VARCHAR(3))
LANGUAGE SQL AS $$
    -- set plan points and deactivate
    UPDATE plans SET updated_date = NOW(), active = false, points = plan_points WHERE id = plan_id;
    -- reset battle active_plan_id
    UPDATE battles SET updated_date = NOW(), active_plan_id = null WHERE id = battle_id;
$$;

-- Delete Battle --
DROP PROCEDURE IF EXISTS delete_battle(battle_id UUID);
CREATE PROCEDURE delete_battle(battle_id UUID)
LANGUAGE SQL AS $$
    DELETE FROM plans WHERE battle_id = battle_id;
    DELETE FROM battles_warriors WHERE battle_id = battle_id;
    DELETE FROM battles WHERE id = battle_id;
$$;