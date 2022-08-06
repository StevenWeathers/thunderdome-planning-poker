DROP FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode TEXT, facilitatorCode TEXT);
CREATE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode TEXT, facilitatorCode TEXT) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code, facilitator_code)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode) RETURNING id INTO storyId;
    INSERT INTO storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION create_retro(
    userId UUID,
    retroName VARCHAR(256),
    fmt VARCHAR(32),
    joinCode TEXT,
    facilitatorCode TEXT,
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
);
CREATE FUNCTION create_retro(
    userId UUID,
    retroName VARCHAR(256),
    fmt VARCHAR(32),
    joinCode TEXT,
    facilitatorCode TEXT,
    maxVotes SMALLINT,
    brainstormVisibility VARCHAR(12)
) RETURNS UUID
AS $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO retro_user (retro_id, user_id) VALUES (retroId, userId);

    RETURN retroId;
END;
$$ LANGUAGE plpgsql;