-- Create a Storyboard
DROP FUNCTION create_storyboard(UUID, VARCHAR);
CREATE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code) VALUES (ownerId, storyboardName, joinCode) RETURNING id INTO storyId;

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

-- Edit a Storyboard
CREATE PROCEDURE edit_storyboard(storyboardId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard SET name = storyboardName, join_code = joinCode, updated_date = NOW()
        WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Edit a Storyboard Story comment --
CREATE PROCEDURE story_comment_edit(commentId UUID, comment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard_story_comment SET comment = comment
        WHERE id = commentId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Delete a Storyboard Story comment --
CREATE PROCEDURE story_comment_delete(commentId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_story_comment WHERE id = commentId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;