-- Demote Battle Leader --
CREATE OR REPLACE PROCEDURE demote_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles_leaders WHERE battle_id = battleId AND user_id = leaderId;
    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;