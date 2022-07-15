ALTER TABLE retro DROP COLUMN facilitator_code;
ALTER TABLE retro ALTER COLUMN join_code TYPE VARCHAR(128);
ALTER TABLE storyboard DROP COLUMN facilitator_code;
ALTER TABLE storyboard ALTER COLUMN join_code TYPE VARCHAR(128);
ALTER TABLE battles ALTER COLUMN join_code TYPE VARCHAR(128);
ALTER TABLE battles ALTER COLUMN leader_code TYPE VARCHAR(128);

DROP FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    IN joinCode TEXT,
    IN leaderCode TEXT,
    OUT battleId UUID
)
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

DROP FUNCTION create_retro(
    userId UUID,
    retroName VARCHAR(256),
    fmt VARCHAR(32),
    joinCode TEXT,
    facilitatorCode TEXT,
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
);
CREATE FUNCTION create_retro(
    userId UUID,
    retroName VARCHAR(256),
    fmt VARCHAR(32),
    joinCode VARCHAR(128),
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
) RETURNS UUID
AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO retro_facilitator (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE edit_retro(
    retroId UUID,
    retroName VARCHAR(256),
    joinCode TEXT,
    facilitatorCode TEXT,
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
);
CREATE PROCEDURE edit_retro(
    retroId UUID,
    retroName VARCHAR(256),
    joinCode VARCHAR(128),
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro
    SET name = retroName, join_code = joinCode, max_votes = maxVotes,
        brainstorm_visibility = brainstormVisibility, updated_date = NOW()
    WHERE id = retroId;

    COMMIT;
END;
$$;

DROP FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode TEXT, facilitatorCode TEXT);
CREATE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code)
        VALUES (ownerId, storyboardName, joinCode) RETURNING id INTO storyId;
    INSERT INTO storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE edit_storyboard(storyboardId UUID, storyboardName VARCHAR(256), joinCode TEXT, facilitatorCode TEXT);
CREATE PROCEDURE edit_storyboard(storyboardId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard SET name = storyboardName, join_code = joinCode, updated_date = NOW()
        WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Get Retros by User ID
DROP FUNCTION get_retros_by_user(userId UUID);
CREATE FUNCTION get_retros_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID, format VARCHAR(32), phase VARCHAR(16), join_code VARCHAR(128), created_date timestamptz, updated_date timestamptz
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id, b.format, b.phase, b.join_code, b.created_date, b.updated_date
		FROM retro b
		LEFT JOIN retro_user su ON b.id = su.retro_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;