ALTER TABLE storyboard_story
   DROP CONSTRAINT storyboard_story_column_id_sort_order_key
 , ADD  CONSTRAINT storyboard_story_column_id_sort_order_key UNIQUE(column_id, sort_order);

-- Move a Storyboard Story to a new column and/or goal --
CREATE OR REPLACE PROCEDURE move_story(storyId UUID, goalId UUID, columnId UUID, placeBefore TEXT)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
DECLARE srcGoalId UUID;
DECLARE srcColumnId UUID;
DECLARE srcSortOrder INTEGER;
DECLARE storyName VARCHAR(256);
DECLARE storyColor VARCHAR(32);
DECLARE storyContent TEXT;
DECLARE createdDate TIMESTAMP;
DECLARE targetSortOrder INTEGER;
BEGIN
    -- Get Story current details
    SELECT
        storyboard_id, goal_id, column_id, sort_order, name, color, content, created_date
    INTO
        storyboardId, srcGoalId, srcColumnId, srcSortOrder, storyName, storyColor, storyContent, createdDate
    FROM storyboard_story WHERE id = storyId;

    -- Get target sort order
    IF placeBefore = '' THEN
        SELECT coalesce(max(sort_order), 0) + 1 INTO targetSortOrder FROM storyboard_story WHERE column_id = columnId;
    ELSE
        SELECT sort_order INTO targetSortOrder FROM storyboard_story WHERE column_id = columnId AND id = placeBefore::UUID;
    END IF;

    -- Remove from source column
    DELETE FROM storyboard_story WHERE id = storyId;
    -- Update sort order in src column
    UPDATE storyboard_story ss SET sort_order = (t.sort_order - 1)
    FROM (
        SELECT id, sort_order FROM storyboard_story
        WHERE column_id = srcColumnId AND sort_order > srcSortOrder
        ORDER BY sort_order ASC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Update sort order for any story that should come after newly moved story
    UPDATE storyboard_story ss SET sort_order = (t.sort_order + 1)
    FROM (
        SELECT id, sort_order FROM storyboard_story
        WHERE column_id = columnId AND sort_order >= targetSortOrder
        ORDER BY sort_order DESC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Finally, insert new story in its ordered place
    INSERT INTO
        storyboard_story (
            storyboard_id, goal_id, column_id, sort_order, name, color, content, created_date
        )
    VALUES (
        storyBoardId, goalId, columnId, targetSortOrder, storyName, storyColor, storyContent, createdDate
    );

    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;