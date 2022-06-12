-- Create a Storyboard
DROP FUNCTION create_storyboard(UUID, VARCHAR, VARCHAR);
CREATE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name) VALUES (ownerId, storyboardName) RETURNING id INTO storyId;

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE edit_storyboard(UUID, VARCHAR, VARCHAR);
DROP PROCEDURE story_comment_edit(UUID, TEXT);
DROP PROCEDURE story_comment_delete(UUID);