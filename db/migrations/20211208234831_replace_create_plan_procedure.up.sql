DROP PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT);


CREATE OR REPLACE PROCEDURE create_plan(battleId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO plans (battle_id, name, type, reference_id, link, description, acceptance_criteria)
    VALUES (battleId, planName, planType, referenceId, planLink, planDescription, acceptanceCriteria);

    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;
