-- Set Retro Owner --
CREATE OR REPLACE PROCEDURE set_retro_owner(retroId UUID, ownerId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro SET updated_date = NOW(), owner_id = ownerId WHERE id = retroId;
END;
$$;

-- Set Retro Phase --
CREATE OR REPLACE PROCEDURE set_retro_phase(retroId UUID, nextPhase VARCHAR(16))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro SET updated_date = NOW(), phase = nextPhase WHERE id = retroId;
END;
$$;

-- Delete Retro --
CREATE OR REPLACE PROCEDURE delete_retro(retroId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retro WHERE id = retroId;

    COMMIT;
END;
$$;

-- Create a Retro
CREATE FUNCTION create_retro(ownerId UUID, retroName VARCHAR(256), format VARCHAR(32), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code) VALUES (ownerId, retroName, format, joinCode) RETURNING id INTO retroId;

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;

-- Get Retros by User ID
CREATE FUNCTION get_retros_by_user(userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), owner_id UUID, format VARCHAR(32), phase VARCHAR(16), join_code VARCHAR(128)
) AS $$
BEGIN
    RETURN QUERY
        SELECT b.id, b.name, b.owner_id, b.format, b.phase, b.join_code
		FROM retro b
		LEFT JOIN retro_user su ON b.id = su.retro_id WHERE su.user_id = userId AND su.abandoned = false
		GROUP BY b.id ORDER BY b.created_date DESC;
END;
$$ LANGUAGE plpgsql;

-- Get Retro Users
CREATE FUNCTION get_retro_users(retroId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL, avatar VARCHAR(128), email VARCHAR(320)
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			u.id, u.name, su.active, u.avatar, COALESCE(u.email, '')
		FROM retro_user su
		LEFT JOIN users u ON su.user_id = u.id
		WHERE su.retro_id = retroId
		ORDER BY u.name;
END;
$$ LANGUAGE plpgsql;

-- Get Retro User by id
CREATE FUNCTION get_retro_user(retroId UUID, userId UUID) RETURNS table (
    id UUID, name VARCHAR(256), active BOOL
) AS $$
BEGIN
    RETURN QUERY
        SELECT
			w.id, w.name, coalesce(su.active, FALSE)
		FROM users w
		LEFT JOIN retro_user su ON su.user_id = w.id AND su.retro_id = retroId
		WHERE w.id = userId;
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

-- Add Retro to Team --
CREATE OR REPLACE FUNCTION team_retro_add(
    IN teamId UUID,
    IN retroId UUID
) RETURNS void AS $$
BEGIN
    INSERT INTO team_retro (team_id, retro_id) VALUES (teamId, retroId);
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Remove Retro from Team --
CREATE OR REPLACE FUNCTION team_retro_remove(
    IN teamId UUID,
    IN retroId UUID
) RETURNS void AS $$
BEGIN
    DELETE FROM team_retro WHERE retro_id = retroId AND team_id = teamId;
    UPDATE team SET updated_date = NOW() WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

-- Clean up Retros older than X Days --
CREATE OR REPLACE PROCEDURE clean_retros(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retro WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;

CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_users SET active = false WHERE active = true;
    UPDATE retro_user SET active = false WHERE active = true;
END;
$$;