-- Edit a Storyboard Story comment --
DROP PROCEDURE story_comment_edit(storyboardId UUID, commentId UUID, comment TEXT);
CREATE PROCEDURE story_comment_edit(storyboardId UUID, commentId UUID, comment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard_story_comment SET comment = comment
        WHERE id = commentId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;
