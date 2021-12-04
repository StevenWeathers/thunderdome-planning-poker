CREATE OR REPLACE FUNCTION user_session_get(
    IN sessionId VARCHAR(64)
) RETURNS table(
    id uuid,
    name VARCHAR(64),
    email VARCHAR(320),
    type VARCHAR(128),
    avatar VARCHAR(128),
    verified BOOLEAN,
    notifications_enabled BOOLEAN,
    country VARCHAR(2),
    locale VARCHAR(2),
    company VARCHAR(256),
    job_title VARCHAR(128),
    created_date TIMESTAMP,
    updated_date TIMESTAMP,
    last_active TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY SELECT
        u.id,
        u.name,
        u.email,
        u.type,
        u.avatar,
        u.verified,
        u.notifications_enabled,
        COALESCE(u.country, ''),
        COALESCE(u.locale, ''),
        COALESCE(u.company, ''),
        COALESCE(u.job_title, ''),
        u.created_date,
        u.updated_date,
        u.last_active
    FROM user_session us
    LEFT JOIN users u ON u.id = us.user_id
    WHERE us.session_id = $1 AND NOW() < us.expire_date;
END;
$$ LANGUAGE plpgsql;

CREATE FUNCTION prune_user_sessions()
RETURNS trigger AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM user_session WHERE expire_date < NOW(); -- clean up old sessions
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM user_session', row_count;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prune_user_sessions AFTER INSERT ON user_session EXECUTE PROCEDURE prune_user_sessions();