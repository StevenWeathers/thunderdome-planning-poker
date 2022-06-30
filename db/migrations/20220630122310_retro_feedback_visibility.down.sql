ALTER TABLE retro DROP COLUMN brainstorm_visibility;
ALTER TABLE retro DROP COLUMN max_votes;
DROP TABLE retro_facilitator;

DROP FUNCTION create_retro(UUID, VARCHAR(256), VARCHAR(32), VARCHAR(128), SMALLINT, VARCHAR(12));
CREATE FUNCTION create_retro(ownerId UUID, retroName VARCHAR(256), fmt VARCHAR(32), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code)
    VALUES (ownerId, retroName, fmt, joinCode) RETURNING id INTO retroId;

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE edit_retro(UUID, VARCHAR(256), VARCHAR(128), SMALLINT, VARCHAR(12));
CREATE PROCEDURE edit_retro(retroId UUID, retroName VARCHAR(256), joinCode VARCHAR(128))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro SET name = retroName, join_code = joinCode, updated_date = NOW()
        WHERE id = retroId;

    COMMIT;
END;
$$;

DROP PROCEDURE retro_add_facilitator(UUID, UUID);
DROP PROCEDURE retro_remove_facilitator(UUID, UUID);

-- Set Retro Owner --
CREATE PROCEDURE set_retro_owner(retroId UUID, ownerId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro SET updated_date = NOW(), owner_id = ownerId WHERE id = retroId;
END;
$$;