ALTER TABLE battles ADD COLUMN hide_voter_identity BOOL DEFAULT false;

DROP FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    IN joinCode TEXT,
    IN leaderCode TEXT,
    OUT battleId UUID
);
CREATE FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    IN hideVoterIdentity BOOL,
    IN joinCode TEXT,
    IN leaderCode TEXT,
    OUT battleId UUID
) AS $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code)
        VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding, hideVoterIdentity, joinCode, leaderCode)
        RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO battles_users (battle_id, user_id) VALUES (battleId, leaderId);
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION team_create_battle(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, joincode text, leadercode text, OUT battleid uuid);
CREATE FUNCTION team_create_battle(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity bool, joincode text, leadercode text, OUT battleid uuid) RETURNS uuid
	language plpgsql
as $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code)
        VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding, hidevoteridentity, joinCode, leaderCode)
        RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO battles_users (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO team_battle (team_id, battle_id) VALUES (teamid, battleId);
END;
$$;