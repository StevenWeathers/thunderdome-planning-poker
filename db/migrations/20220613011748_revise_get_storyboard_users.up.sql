-- Get Storyboard Users
DROP FUNCTION get_storyboard_users(storyboardId UUID);
CREATE FUNCTION get_storyboard_users(storyboardId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL, avatar varchar(128), email varchar(320)
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, su.active, w.avatar, COALESCE(w.email, '')
		FROM storyboard_user su
		LEFT JOIN users w ON su.user_id = w.id
		WHERE su.storyboard_id = storyboardId
		ORDER BY w.name;
END;
$$ LANGUAGE plpgsql;