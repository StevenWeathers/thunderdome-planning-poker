-- Create Battle --
DROP FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    OUT battleId UUID
);
CREATE FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    IN joinCode VARCHAR(128),
    IN leaderCode VARCHAR(128),
    OUT battleId UUID
) AS $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, join_code, leader_code)
        VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding, joinCode, leaderCode)
        RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO battles_users (battle_id, user_id) VALUES (battleId, leaderId);
END;
$$ LANGUAGE plpgsql;