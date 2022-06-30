ALTER TABLE retro ADD COLUMN brainstorm_visibility VARCHAR(12) DEFAULT 'visible'; -- visible, concealed, hidden
ALTER TABLE retro ADD COLUMN max_votes SMALLINT DEFAULT 3;

CREATE TABLE "retro_facilitator" (
    "retro_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "created_date" TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ("retro_id","user_id")
);

-- back fill the retro_facilitator table with owner ids
BEGIN;
    INSERT INTO retro_facilitator (retro_id, user_id)
    SELECT id, owner_id FROM retro;
COMMIT;

-- Create a Retro
DROP FUNCTION create_retro(UUID, VARCHAR(256), VARCHAR(32), VARCHAR(128));
CREATE FUNCTION create_retro(
    userId UUID,
    retroName VARCHAR(256),
    fmt VARCHAR(32),
    joinCode VARCHAR(128),
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
) RETURNS UUID
AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO retro_facilitator (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;

-- Edit a Retro
DROP PROCEDURE edit_retro(UUID, VARCHAR(256), VARCHAR(128));
CREATE PROCEDURE edit_retro(
    retroId UUID,
    retroName VARCHAR(256),
    joinCode VARCHAR(128),
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro
    SET name = retroName, join_code = joinCode, max_votes = maxVotes,
        brainstorm_visibility = brainstormVisibility, updated_date = NOW()
    WHERE id = retroId;

    COMMIT;
END;
$$;

CREATE PROCEDURE retro_add_facilitator(retroId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;
END;
$$;

CREATE PROCEDURE retro_remove_facilitator(retroId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retro_facilitator WHERE retro_id = retroId AND user_id = userId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;
END;
$$;

DROP PROCEDURE set_retro_owner(UUID, UUID);