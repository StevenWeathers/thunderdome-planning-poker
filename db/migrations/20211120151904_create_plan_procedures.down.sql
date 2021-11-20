DROP PROCEDURE activate_plan_voting(battleId UUID, planId UUID);
DROP PROCEDURE create_plan(battleId UUID, planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT);
DROP PROCEDURE delete_plan(battleId UUID, planId UUID);
DROP PROCEDURE end_plan_voting(battleId UUID, planId UUID);
DROP PROCEDURE finalize_plan(battleId UUID, planId UUID, planPoints VARCHAR(3));
DROP PROCEDURE revise_plan(planId UUID, planName VARCHAR(256), planType VARCHAR(64), referenceId VARCHAR(128), planLink TEXT, planDescription TEXT, acceptanceCriteria TEXT);
DROP PROCEDURE skip_plan_voting(battleId UUID, planId UUID);