-- Edit a Storyboard Story comment --
DROP PROCEDURE story_comment_edit(commentId UUID, comment TEXT);
CREATE PROCEDURE story_comment_edit(storyboardId UUID, commentId UUID, comment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard_story_comment SET comment = comment
        WHERE id = commentId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Delete a Storyboard Story comment --
DROP PROCEDURE story_comment_delete(commentId UUID);
CREATE PROCEDURE story_comment_delete(storyboardId UUID, commentId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_story_comment WHERE id = commentId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;