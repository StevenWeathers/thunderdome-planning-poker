-- Get API Keys --
CREATE OR REPLACE FUNCTION apikeys_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id text, name VARCHAR(256), email VARCHAR(320), active BOOLEAN, created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT apk.id, apk.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM api_keys apk
		LEFT JOIN users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Insert a new user password reset
CREATE OR REPLACE FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
)
AS $$
BEGIN
    SELECT id, name INTO userId, userName FROM users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Register a new user
CREATE OR REPLACE FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    INSERT INTO users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Register a new user from existing guest
CREATE OR REPLACE FUNCTION register_existing_user(
    IN activeUserId UUID,
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    UPDATE users
    SET
        name = userName,
        email = userEmail,
        password = hashedPassword,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Get Registered Users list --
CREATE OR REPLACE FUNCTION registered_users_list(
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
CREATE OR REPLACE FUNCTION registered_users_email_search(
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

-- Insert a new user apikey
CREATE OR REPLACE FUNCTION user_apikey_add(
    IN apikeyId text,
    IN keyName VARCHAR(256),
    IN userId uuid,
    OUT createdDate timestamp
)
AS $$
BEGIN
    INSERT INTO api_keys (id, name, user_id) VALUES (apikeyId, keyName, userId) RETURNING created_date INTO createdDate;
    UPDATE users SET last_active = NOW() WHERE id = userId;
END;
$$ LANGUAGE plpgsql;

-- Deletes a user apikey
CREATE OR REPLACE PROCEDURE user_apikey_delete(
    apikeyId text,
    userId uuid
)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM api_keys WHERE id = apikeyId AND user_id = userId;
    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Updates a user apikey
CREATE OR REPLACE PROCEDURE user_apikey_update(
    apikeyId text,
    userId uuid,
    keyActive boolean
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE api_keys SET active = keyActive, updated_date = NOW() WHERE id = apikeyId AND user_id = userId;
    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;