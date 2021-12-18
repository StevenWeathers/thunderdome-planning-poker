-- Get Team Users --
DROP FUNCTION team_user_list(uuid,integer,integer);
CREATE OR REPLACE FUNCTION team_user_list(
    IN teamId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), tu.role
        FROM team_user tu
        LEFT JOIN users u ON tu.user_id = u.id
        WHERE tu.team_id = teamId
        ORDER BY tu.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;
-- Get Organization Users --
DROP FUNCTION organization_user_list(uuid,integer,integer);
CREATE OR REPLACE FUNCTION organization_user_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), ou.role
        FROM organization_user ou
        LEFT JOIN users u ON ou.user_id = u.id
        WHERE ou.organization_id = orgId
        ORDER BY ou.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;
-- Get Department Users --
DROP FUNCTION department_user_list(uuid,integer,integer);
CREATE OR REPLACE FUNCTION department_user_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), email VARCHAR(256), role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), du.role
        FROM department_user du
        LEFT JOIN users u ON du.user_id = u.id
        WHERE du.department_id = departmentId
        ORDER BY du.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;