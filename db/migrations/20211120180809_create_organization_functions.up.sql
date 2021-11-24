-- Get Organization --
CREATE OR REPLACE FUNCTION organization_get_by_id(
    IN orgId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        WHERE o.id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization User Role --
CREATE OR REPLACE FUNCTION organization_get_user_role(
    IN userId UUID,
    IN orgId UUID,
    OUT role VARCHAR(16)
) AS $$
BEGIN
    SELECT ou.role INTO role
    FROM organization_user ou
    WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations --
CREATE OR REPLACE FUNCTION organization_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM organization o
        ORDER BY o.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Organizations by User --
CREATE OR REPLACE FUNCTION organization_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP, role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date, ou.role
        FROM organization_user ou
        LEFT JOIN organization o ON ou.organization_id = o.id
        WHERE ou.user_id = userId
        ORDER BY created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization --
DROP FUNCTION IF EXISTS organization_create(IN userId UUID, IN orgName VARCHAR(256), OUT organizationId UUID);
CREATE OR REPLACE FUNCTION organization_create(
    IN userId UUID,
    IN orgName VARCHAR(256)
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
    DECLARE organizationId uuid;
BEGIN
    INSERT INTO organization (name) VALUES (orgName) RETURNING organization.id INTO organizationId;
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
    RETURN QUERY SELECT o.id, o.name, o.created_date, o.updated_date FROM organization o WHERE o.id = organizationID;
END;
$$ LANGUAGE plpgsql;

-- Add User to Organization --
CREATE OR REPLACE FUNCTION organization_user_add(
    IN orgId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (orgId, userId, userRole);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Organization --
CREATE OR REPLACE PROCEDURE organization_user_remove(orgId UUID, userId UUID)
AS $$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM organization_department WHERE organization_id = orgId
    LOOP
        CALL department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM team_user tu WHERE tu.team_id IN (
        SELECT ot.team_id
        FROM organization_team ot
        WHERE ot.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM organization_user WHERE organization_id = orgId AND user_id = userId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Users --
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

-- Get Organization Teams --
CREATE OR REPLACE FUNCTION organization_team_list(
    IN orgId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM organization_team ot
        LEFT JOIN team t ON ot.team_id = t.id
        WHERE ot.organization_id = orgId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Organization Team --
DROP FUNCTION IF EXISTS organization_team_create(IN orgId UUID, IN teamName VARCHAR(256), OUT teamId UUID);
CREATE OR REPLACE FUNCTION organization_team_create(
    IN orgId UUID,
    IN teamName VARCHAR(256)
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
    DECLARE teamId uuid;
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING team.id INTO teamId;
    INSERT INTO organization_team (organization_id, team_id) VALUES (orgId, teamId);
    UPDATE organization SET updated_date = NOW() WHERE organization.id = orgId;
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date FROM team t WHERE t.id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Organization Team User Role --
CREATE OR REPLACE FUNCTION organization_team_user_role(
    IN userId UUID,
    IN orgId UUID,
    IN teamId UUID
) RETURNS table (
    orgRole VARCHAR(16), teamRole VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT ou.role AS orgRole, COALESCE(tu.role, '') AS teamRole
        FROM organization_user ou
        LEFT JOIN team_user tu ON tu.user_id = userId AND tu.team_id = teamId
        WHERE ou.organization_id = orgId AND ou.user_id = userId;
END;
$$ LANGUAGE plpgsql;