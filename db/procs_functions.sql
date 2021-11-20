--
-- Stored Functions
--

-- Get Application Stats e.g. total user and battle counts
CREATE OR REPLACE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER,
    OUT apikey_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO battle_count FROM battles;
    SELECT COUNT(*) INTO plan_count FROM plans;
    SELECT COUNT(*) INTO organization_count FROM organization;
    SELECT COUNT(*) INTO department_count FROM organization_department;
    SELECT COUNT(*) INTO team_count FROM team;
    SELECT COUNT(*) INTO apikey_count FROM api_keys;
END;
$$ LANGUAGE plpgsql;

-- Insert a new user password reset
CREATE OR REPLACE FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
)
AS $$
BEGIN
    SELECT id, name INTO userId, userName FROM users WHERE email = userEmail;
    IF FOUND THEN
        INSERT INTO user_reset (user_id) VALUES (userId) RETURNING reset_id INTO resetId;
    ELSE
        RAISE EXCEPTION 'Nonexistent User --> %', userEmail USING HINT = 'Please check your Email';
    END IF;
END;
$$ LANGUAGE plpgsql;

-- Register a new user
CREATE OR REPLACE FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    INSERT INTO users (name, email, password, type)
    VALUES (userName, userEmail, hashedPassword, userType)
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Register a new user from existing guest
CREATE OR REPLACE FUNCTION register_existing_user(
    IN activeUserId UUID,
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
)
AS $$
BEGIN
    UPDATE users
    SET
        name = userName,
        email = userEmail,
        password = hashedPassword,
        type = userType,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = activeUserId
    RETURNING id INTO userId;

    INSERT INTO user_verify (user_id) VALUES (userId) RETURNING verify_id INTO verifyId;
END;
$$ LANGUAGE plpgsql;

-- Create Battle --
CREATE OR REPLACE FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    OUT battleId UUID
) AS $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding) VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding) RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
END;
$$ LANGUAGE plpgsql;

-- Add Battle Leaders by Emails --
CREATE OR REPLACE FUNCTION add_battle_leaders_by_email(
    IN battleId UUID,
    IN leaderEmails TEXT,
    OUT leaders JSONB
) AS $$
DECLARE
    emails TEXT[];
    leaderEmail TEXT;
BEGIN
    select into emails regexp_split_to_array(leaderEmails,',');
    FOREACH leaderEmail IN ARRAY emails
    LOOP
        INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, (
            SELECT id FROM users WHERE email = leaderEmail
        ));
    END LOOP;

    SELECT CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END
    FROM battles_leaders bl WHERE bl.battle_id = battleId INTO leaders;
END;
$$ LANGUAGE plpgsql;

-- Get a list of countries
CREATE OR REPLACE FUNCTION countries_active() RETURNS table (
    country VARCHAR(2)
) AS $$
BEGIN
    RETURN QUERY SELECT ac.country FROM active_countries ac;
END;
$$ LANGUAGE plpgsql;

-- Get API Keys --
CREATE OR REPLACE FUNCTION apikeys_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id text, name VARCHAR(256), email VARCHAR(320), active BOOLEAN, created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT apk.id, apk.name, u.email, apk.active, apk.created_date, apk.updated_date
		FROM api_keys apk
		LEFT JOIN users u ON apk.user_id = u.id
		ORDER BY apk.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Registered Users list --
CREATE OR REPLACE FUNCTION registered_users_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid, name VARCHAR(64), email VARCHAR(320), type VARCHAR(128), avatar VARCHAR(128), verified BOOLEAN, country VARCHAR(2), company VARCHAR(256), job_title VARCHAR(128)
) AS $$
BEGIN
    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, '')
		FROM users u
		WHERE u.email IS NOT NULL
		ORDER BY u.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Search Registered Users for those like email --
CREATE OR REPLACE FUNCTION registered_users_email_search(
    IN email_search VARCHAR(320),
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table(
    id uuid,
    name VARCHAR(64),
    email VARCHAR(320),
    type VARCHAR(128),
    avatar VARCHAR(128),
    verified BOOLEAN,
    country VARCHAR(2),
    company VARCHAR(256),
    job_title VARCHAR(128),
    count INTEGER
) AS $$
    DECLARE count INTEGER;
BEGIN
    SELECT count(*)
    FROM users u
    WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%') INTO count;

    RETURN QUERY
        SELECT u.id, u.name, COALESCE(u.email, ''), u.type, u.avatar, u.verified, COALESCE(u.country, ''), COALESCE(u.company, ''), COALESCE(u.job_title, ''), count
        FROM users u
        WHERE u.email IS NOT NULL AND u.email LIKE ('%' || email_search || '%')
        ORDER BY u.created_date
        LIMIT l_limit
        OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

--
-- ORGANIZATIONS --
--

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
CREATE OR REPLACE FUNCTION organization_create(
    IN userId UUID,
    IN orgName VARCHAR(256),
    OUT organizationId UUID
) AS $$
BEGIN
    INSERT INTO organization (name) VALUES (orgName) RETURNING id INTO organizationId;
    INSERT INTO organization_user (organization_id, user_id, role) VALUES (organizationId, userId, 'ADMIN');
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
CREATE OR REPLACE FUNCTION organization_team_create(
    IN orgId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO organization_team (organization_id, team_id) VALUES (orgId, teamId);
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
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

--
-- DEPARTMENTS --
--

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
CREATE OR REPLACE FUNCTION department_create(
    IN orgId UUID,
    IN departmentName VARCHAR(256),
    OUT departmentId UUID
) AS $$
BEGIN
    INSERT INTO organization_department (name, organization_id) VALUES (departmentName, orgId) RETURNING id INTO departmentId;
    UPDATE organization SET updated_date = NOW() WHERE id = orgId;
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
CREATE OR REPLACE FUNCTION department_team_create(
    IN departmentId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO department_team (department_id, team_id) VALUES (departmentId, teamId);
    UPDATE organization_department SET updated_date = NOW() WHERE id = departmentId;
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

--
-- TEAMS --
--

-- Get Team --
CREATE OR REPLACE FUNCTION team_get_by_id(
    IN teamId UUID
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT o.id, o.name, o.created_date, o.updated_date
        FROM team o
        WHERE o.id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Get Team User Role --
CREATE OR REPLACE FUNCTION team_get_user_role(
    IN userId UUID,
    IN teamId UUID
) RETURNS table (
    role VARCHAR(16)
) AS $$
BEGIN
    RETURN QUERY
        SELECT tu.role
        FROM team_user tu
        WHERE tu.team_id = teamId AND tu.user_id = userId;
END;
$$ LANGUAGE plpgsql;

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
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Get Teams by User --
CREATE OR REPLACE FUNCTION team_list_by_user(
    IN userId UUID,
    IN l_limit INTEGER,
    IN l_offset INTEGER
) RETURNS table (
    id UUID, name VARCHAR(256), created_date TIMESTAMP, updated_date TIMESTAMP
) AS $$
BEGIN
    RETURN QUERY
        SELECT t.id, t.name, t.created_date, t.updated_date
        FROM team_user tu
        LEFT JOIN team t ON tu.team_id = t.id
        WHERE tu.user_id = userId
        ORDER BY t.created_date
		LIMIT l_limit
		OFFSET l_offset;
END;
$$ LANGUAGE plpgsql;

-- Create Team --
CREATE OR REPLACE FUNCTION team_create(
    IN userId UUID,
    IN teamName VARCHAR(256),
    OUT teamId UUID
) AS $$
BEGIN
    INSERT INTO team (name) VALUES (teamName) RETURNING id INTO teamId;
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, 'ADMIN');
END;
$$ LANGUAGE plpgsql;

-- Get Team Users --
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

-- Add User to Team --
CREATE OR REPLACE FUNCTION team_user_add(
    IN teamId UUID,
    IN userId UUID,
    IN userRole VARCHAR(16)
) RETURNS void AS $$
BEGIN
    INSERT INTO team_user (team_id, user_id, role) VALUES (teamId, userId, userRole);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove User from Team --
CREATE OR REPLACE PROCEDURE team_user_remove(teamId UUID, userId UUID)
AS $$
BEGIN
    DELETE FROM team_user WHERE team_id = teamId AND user_id = userId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

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

-- Add Battle to Team --
CREATE OR REPLACE FUNCTION team_battle_add(
    IN teamId UUID,
    IN battleId UUID
) RETURNS void AS $$
BEGIN
    INSERT INTO team_battle (team_id, battle_id) VALUES (teamId, battleId);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove Battle from Team --
CREATE OR REPLACE FUNCTION team_battle_remove(
    IN teamId UUID,
    IN battleId UUID
) RETURNS void AS $$
BEGIN
    DELETE FROM team_battle WHERE battle_id = battleId AND team_id = teamId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Delete Team --
CREATE OR REPLACE PROCEDURE team_delete(teamId UUID)
AS $$
BEGIN
    DELETE FROM team WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;