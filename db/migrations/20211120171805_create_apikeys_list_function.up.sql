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