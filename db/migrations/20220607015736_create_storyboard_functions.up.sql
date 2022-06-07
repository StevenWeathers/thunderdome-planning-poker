-- Create a Storyboard
CREATE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name) VALUES (ownerId, storyboardName) RETURNING id INTO storyId;

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

-- Get Storyboards by User ID
CREATE FUNCTION get_storyboards_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id
		FROM storyboard b
		LEFT JOIN storyboard_user su ON b.id = su.storyboard_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;

-- Get a Storyboards Goals --
CREATE FUNCTION get_storyboard_goals(storyboardId UUID) RETURNS table (
    id UUID, sort_order INTEGER, name VARCHAR(256), columns JSON
) AS $$
BEGIN
    RETURN QUERY
        SELECT
            sg.id,
            sg.sort_order,
            sg.name,
            COALESCE(json_agg(to_jsonb(t) - 'goal_id' ORDER BY t.sort_order) FILTER (WHERE t.id IS NOT NULL), '[]') AS columns
        FROM storyboard_goal sg
        LEFT JOIN (
            SELECT
                sc.*,
                COALESCE(
                    json_agg(stss ORDER BY stss.sort_order) FILTER (WHERE stss.id IS NOT NULL), '[]'
                ) AS stories
            FROM storyboard_column sc
            LEFT JOIN (
                SELECT
                    ss.*,
                    COALESCE(
                        json_agg(stcm ORDER BY stcm.created_date) FILTER (WHERE stcm.id IS NOT NULL), '[]'
                    ) AS comments
                FROM storyboard_story ss
                LEFT JOIN storyboard_story_comment stcm ON stcm.story_id = ss.id
                GROUP BY ss.id
            ) stss ON stss.column_id = sc.id
            GROUP BY sc.id
        ) t ON t.goal_id = sg.id
        WHERE sg.storyboard_id = storyboardId
        GROUP BY sg.id
        ORDER BY sg.sort_order;
END;
$$ LANGUAGE plpgsql;

-- Get Storyboard Users
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

-- Get Storyboard Personas
CREATE FUNCTION get_storyboard_personas(storyboardId UUID) RETURNS table (
    id UUID,
    name VARCHAR(256),
    role VARCHAR(256),
    description TEXT
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			p.id, p.name, p.role, p.description
		FROM storyboard_persona p
		WHERE p.storyboard_id = storyboardId;
END;
$$ LANGUAGE plpgsql;

-- Get Storyboard User by id
CREATE FUNCTION get_storyboard_user(storyboardId UUID, userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, coalesce(su.active, FALSE)
		FROM users w
		LEFT JOIN storyboard_user su ON su.user_id = w.id AND su.storyboard_id = storyboardId
		WHERE w.id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Team Storyboards --
CREATE FUNCTION team_storyboard_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name
        FROM team_storyboard tb
        LEFT JOIN storyboard b ON tb.storyboard_id = b.id
        WHERE tb.team_id = teamId
        ORDER BY tb.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Add Storyboard to Team --
CREATE FUNCTION team_storyboard_add(
    IN teamId UUID,
    IN storyboardId UUID
) RETURNS void AS $$
BEGIN
    INSERT INTO team_storyboard (team_id, storyboard_id) VALUES (teamId, storyboardId);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove Storyboard from Team --
CREATE FUNCTION team_storyboard_remove(
    IN teamId UUID,
    IN storyboardId UUID
) RETURNS void AS $$
BEGIN
    DELETE FROM team_storyboard WHERE storyboard_id = storyboardId AND team_id = teamId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;