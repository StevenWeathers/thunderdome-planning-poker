-- Set Battle Leader --
CREATE OR REPLACE PROCEDURE set_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;