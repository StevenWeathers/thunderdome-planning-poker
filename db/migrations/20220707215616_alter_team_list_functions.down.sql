-- Get Team Battles --
CREATE OR REPLACE FUNCTION team_battle_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name
        FROM team_battle tb
        LEFT JOIN battles b ON tb.battle_id = b.id
        WHERE tb.team_id = teamId
        ORDER BY tb.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Team Retros --
CREATE OR REPLACE FUNCTION team_retro_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), format VARCHAR(32), phase VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.format, b.phase
        FROM team_retro tb
        LEFT JOIN retro b ON tb.retro_id = b.id
        WHERE tb.team_id = teamId
        ORDER BY tb.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Team Storyboards --
CREATE OR REPLACE FUNCTION team_storyboard_list(
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