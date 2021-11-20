DROP PROCEDURE delete_battle(battleId UUID);
DROP PROCEDURE demote_battle_leader(battleId UUID, leaderId UUID);
DROP PROCEDURE retract_user_vote(planId UUID, userId UUID);
DROP PROCEDURE set_battle_leader(battleId UUID, leaderId UUID);
DROP PROCEDURE set_user_vote(planId UUID, userId UUID, userVote VARCHAR(3));