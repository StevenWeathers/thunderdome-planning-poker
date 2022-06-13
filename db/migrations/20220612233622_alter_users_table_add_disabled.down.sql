ALTER TABLE users DROP COLUMN disabled;

DROP PROCEDURE user_disable(userId UUID);
DROP PROCEDURE user_enable(userId UUID);

-- Get Registered Users list --
DROP FUNCTION registered_users_list(INTEGER, INTEGER);
CREATE FUNCTION registered_users_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), avatar VARCHAR(128), verified BOOLEAN, country VARCHAR(2), company VARCHAR(256), job_title VARCHAR(128)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, '')
		FROM users u
		WHERE u.email IS NOT NULL
		ORDER BY u.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Search Registered Users for those like email --
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
    WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%') INTO count;

    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), count
        FROM users u
        WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%')
        ORDER BY u.created_date
        LIMIT l_limit
        OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;