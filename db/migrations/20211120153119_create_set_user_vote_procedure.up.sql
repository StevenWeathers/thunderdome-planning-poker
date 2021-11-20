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

    COMMIT;
END;
$$;