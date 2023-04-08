-- Get Teams --
CREATE OR REPLACE FUNCTION team_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team t
        LEFT JOIN department_team dt on t.id = dt.team_id
        LEFT JOIN organization_team ot on t.id = ot.team_id
        WHERE dt.team_id IS NULL AND ot.team_id IS NULL
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Teams --
CREATE OR REPLACE FUNCTION team_list_count(OUT count INTEGER)
AS $$
BEGIN
    SELECT count(t.id) INTO count
    FROM team t
    LEFT JOIN department_team dt on t.id = dt.team_id
    LEFT JOIN organization_team ot on t.id = ot.team_id
    WHERE dt.team_id IS NULL AND ot.team_id IS NULL;
END;
$$ LANGUAGE plpgsql;