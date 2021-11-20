-- Delete Battle --
CREATE OR REPLACE PROCEDURE delete_battle(battleId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles WHERE id = battleId;

    COMMIT;
END;
$$;

-- Demote Battle Leader --
CREATE OR REPLACE PROCEDURE demote_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles_leaders WHERE battle_id = battleId AND user_id = leaderId;
    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Retract User Vote --
CREATE OR REPLACE PROCEDURE retract_user_vote(planId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT coalesce(json_agg(data), '[]'::JSON)
        FROM (
            SELECT coalesce(oldVote."warriorId") AS "warriorId", coalesce(oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::UsersVote,p1.votes) AS oldVote
            WHERE oldVote."warriorId" != userId
        ) data
    )
    WHERE p1.id = planId;

    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Set Battle Leader --
CREATE OR REPLACE PROCEDURE set_battle_leader(battleId UUID, leaderId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    UPDATE battles SET updated_date = NOW() WHERE id = battleId;
END;
$$;

-- Set User Vote --
CREATE OR REPLACE PROCEDURE set_user_vote(planId UUID, userId UUID, userVote VARCHAR(3))
LANGUAGE plpgsql AS $$
BEGIN
	UPDATE plans p1
    SET votes = (
        SELECT json_agg(data)
        FROM (
            SELECT coalesce(newVote."warriorId", oldVote."warriorId") AS "warriorId", coalesce(newVote.vote, oldVote.vote) AS vote
            FROM jsonb_populate_recordset(null::UsersVote,p1.votes) AS oldVote
            FULL JOIN jsonb_populate_recordset(null::UsersVote,
                ('[{"warriorId":"'|| userId::TEXT ||'", "vote":"'|| userVote ||'"}]')::JSONB
            ) AS newVote
            ON newVote."warriorId" = oldVote."warriorId"
        ) data
    )
    WHERE p1.id = planId;

    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;