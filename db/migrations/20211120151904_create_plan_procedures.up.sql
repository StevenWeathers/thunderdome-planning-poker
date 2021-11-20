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

-- Create a Battle Plan --
CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (id, battle_id, name, type, reference_id, link, description, acceptance_criteria)
    VALUES (planId, battleId, planName, planType, referenceId, planLink, planDescription, acceptanceCriteria);

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
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