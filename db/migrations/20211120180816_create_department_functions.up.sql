-- Get Department --
CREATE OR REPLACE FUNCTION department_get_by_id(
    IN departmentId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT od.id, od.name, od.created_date, od.updated_date
        FROM organization_department od
        WHERE od.id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department User Role --
CREATE OR REPLACE FUNCTION department_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Departments --
CREATE OR REPLACE FUNCTION department_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT d.id, d.name, d.created_date, d.updated_date
        FROM organization_department d
        WHERE d.organization_id = orgId
        ORDER BY d.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Department --
DROP FUNCTION IF EXISTS department_create(IN orgId UUID, IN departmentName VARCHAR(256), OUT departmentId UUID);
CREATE OR REPLACE FUNCTION department_create(
    IN orgId UUID,
    IN departmentName VARCHAR(256)
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
    DECLARE departmentId uuid;
BEGIN
    INSERT INTO organization_department (name, organization_id) VALUES (departmentName, orgId) RETURNING organization_department.id INTO departmentId;
    UPDATE organization SET updated_date = NOW() WHERE organization.id = orgId;
    RETURN QUERY SELECT d.id, d.name, d.created_date, d.updated_date FROM organization_department d WHERE d.id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Teams --
CREATE OR REPLACE FUNCTION department_team_list(
    IN departmentId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM department_team dt
        LEFT JOIN team t ON dt.team_id = t.id
        WHERE dt.department_id = departmentId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Department Team --
DROP FUNCTION IF EXISTS department_team_create(IN departmentId UUID, IN teamName VARCHAR(256), OUT teamId UUID);
CREATE OR REPLACE FUNCTION department_team_create(
    IN departmentId UUID,
    IN teamName VARCHAR(256)
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING team.id INTO teamId;
    INSERT INTO department_team (department_id, team_id) VALUES (departmentId, teamId);
    UPDATE organization_department SET updated_date = NOW() WHERE organization_department.id = departmentId;
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM team t WHERE t.id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Team User Role --
CREATE OR REPLACE FUNCTION department_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN departmentId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), departmentRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(du.role, '') AS departmentRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN department_user du ON du.user_id = userId AND du.department_id = departmentId
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Department Users --
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

-- Add User to Department --
CREATE OR REPLACE FUNCTION department_user_add(
    IN departmentId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
DECLARE orgId UUID;
BEGIN
    SELECT organization_id INTO orgId FROM organization_user WHERE user_id = userId;

    IF orgId IS NULL THEN
        RAISE EXCEPTION 'User not in Organization -> %', userId USING HINT = 'Please add user to Organization before department';
    END IF;

    INSERT INTO department_user (department_id, user_id, role) VALUES (departmentId, userId, userRole);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Department --
CREATE OR REPLACE PROCEDURE department_user_remove(departmentId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT dt.team_id
        FROM department_team dt
        WHERE dt.department_id = departmentId
    ) AND tu.user_id = userId;
    DELETE FROM department_user WHERE department_id = departmentId AND user_id = userId;
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;