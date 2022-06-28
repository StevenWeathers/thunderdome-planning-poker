-- Get Storyboards by User ID
DROP FUNCTION get_storyboards_by_user(userId UUID);
CREATE FUNCTION get_storyboards_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID, created_date timestamptz, updated_date timestamptz
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id, b.created_date, b.updated_date
		FROM storyboard b
		LEFT JOIN storyboard_user su ON b.id = su.storyboard_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;

-- Get Retros by User ID
DROP FUNCTION get_retros_by_user(userId UUID);
CREATE FUNCTION get_retros_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID, format VARCHAR(32), phase VARCHAR(16), join_code VARCHAR(128), created_date timestamptz, updated_date timestamptz
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id, b.format, b.phase, b.join_code, b.created_date, b.updated_date
		FROM retro b
		LEFT JOIN retro_user su ON b.id = su.retro_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;