ALTER TABLE users ADD COLUMN disabled bool DEFAULT false;

-- Disable a user
CREATE PROCEDURE user_disable(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET disabled = true, updated_date = NOW()
        WHERE id = userId;
    DELETE FROM user_session WHERE user_id = userId;

    COMMIT;
END;
$$;

-- Enable a user
CREATE PROCEDURE user_enable(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET disabled = false, updated_date = NOW()
        WHERE id = userId;

    COMMIT;
END;
$$;

-- Get Registered Users list --
DROP FUNCTION registered_users_list(INTEGER, INTEGER);
CREATE FUNCTION registered_users_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), avatar VARCHAR(128), verified BOOLEAN, country VARCHAR(2), company VARCHAR(256), job_title VARCHAR(128), disabled bool
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), u.disabled
		FROM users u
		WHERE u.email IS NOT NULL
		ORDER BY u.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

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