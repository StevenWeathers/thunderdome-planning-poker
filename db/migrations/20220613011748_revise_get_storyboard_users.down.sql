-- Get Storyboard Users
DROP FUNCTION get_storyboard_users(storyboardId UUID);
CREATE FUNCTION get_storyboard_users(storyboardId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, su.active
		FROM storyboard_user su
		LEFT JOIN users w ON su.user_id = w.id
		WHERE su.storyboard_id = storyboardId
		ORDER BY w.name;
END;
$$ LANGUAGE plpgsql;