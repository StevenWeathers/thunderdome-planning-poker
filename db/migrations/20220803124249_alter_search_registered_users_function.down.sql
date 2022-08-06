DROP FUNCTION registered_users_email_search(VARCHAR,INTEGER,INTEGER);
CREATE FUNCTION registered_users_email_search(
    IN email_search VARCHAR(320),
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid,
    name VARCHAR(64),
    email VARCHAR(320),
    type VARCHAR(128),
    avatar VARCHAR(128),
    verified BOOLEAN,
    country VARCHAR(2),
    company VARCHAR(256),
    job_title VARCHAR(128),
    count INTEGER
) AS $$
    DECLARE count INTEGER;
BEGIN
    SELECT count(*)
    FROM users u
    WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%') AND u.disabled IS FALSE INTO count;

    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), count
        FROM users u
        WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%') AND u.disabled IS FALSE
        ORDER BY u.created_date
        LIMIT l_limit
        OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION add_battle_leaders_by_email(
    IN battleId UUID,
    IN leaderEmails TEXT,
    OUT leaders JSONB
) AS $$
DECLARE
    emails TEXT[];
    leaderEmail TEXT;
BEGIN
    select into emails regexp_split_to_array(leaderEmails,',');
    FOREACH leaderEmail IN ARRAY emails
    LOOP
        INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, (
            SELECT id FROM users WHERE email = leaderEmail
        ));
    END LOOP;

    SELECT CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END
    FROM battles_leaders bl WHERE bl.battle_id = battleId INTO leaders;
END;
$$ LANGUAGE plpgsql;