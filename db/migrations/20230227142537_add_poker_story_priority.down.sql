DROP PROCEDURE create_plan(battleId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT, planPriority INTEGER);

CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (battle_id, name, type, reference_id, link, description, acceptance_criteria)
    VALUES (battleId, planName, planType, referenceId, planLink, planDescription, acceptanceCriteria);

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Revise Plan --
DROP PROCEDURE revise_plan(planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT, planPriority INTEGER);
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

ALTER TABLE plans DROP COLUMN priority;