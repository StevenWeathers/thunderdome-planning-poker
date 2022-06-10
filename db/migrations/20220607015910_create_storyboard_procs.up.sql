-- Set Storyboard Owner --
CREATE PROCEDURE set_storyboard_owner(storyboardId UUID, ownerId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard SET updated_date = NOW(), owner_id = ownerId WHERE id = storyboardId;
END;
$$;

-- Revise Storyboard ColorLegend --
CREATE PROCEDURE revise_color_legend(storyboardId UUID, colorLegend JSONB)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard SET updated_date = NOW(), color_legend = colorLegend WHERE id = storyboardId;
END;
$$;

-- Delete Storyboard --
CREATE PROCEDURE delete_storyboard(storyboardId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Create a Storyboard Goal --
CREATE PROCEDURE create_storyboard_goal(storyBoardId UUID, goalName VARCHAR(256))
LANGUAGE plpgsql AS $$
DECLARE sortOrder INTEGER;
BEGIN
    sortOrder := (SELECT coalesce(MAX(sort_order), 0) FROM storyboard_goal WHERE storyboard_id = storyBoardId) + 1;
    INSERT INTO
        storyboard_goal
        (storyboard_id, sort_order, name)
        VALUES (storyBoardId, sortOrder, goalName);

    UPDATE storyboard SET updated_date = NOW() WHERE id = storyBoardId;
END;
$$;

-- Revise a Storyboard Goal --
CREATE PROCEDURE update_storyboard_goal(goalId UUID, goalName VARCHAR(256))
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_goal SET name = goalName, updated_date = NOW() WHERE id = goalId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Delete a Storyboard Goal --
CREATE PROCEDURE delete_storyboard_goal(goalId UUID)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
DECLARE sortOrder INTEGER;
BEGIN
    SELECT sort_order, storyboard_id INTO sortOrder, storyboardId FROM storyboard_goal WHERE id = goalId;

    DELETE FROM storyboard_story WHERE goal_id = goalId;
    DELETE FROM storyboard_column WHERE goal_id = goalId;
    DELETE FROM storyboard_goal WHERE id = goalId;
    UPDATE storyboard_goal sg SET sort_order = (sg.sort_order - 1) WHERE sg.storyboard_id = storyBoardId AND sg.sort_order > sortOrder;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Create a Storyboard Column --
CREATE PROCEDURE create_storyboard_column(storyBoardId UUID, goalId UUID)
LANGUAGE plpgsql AS $$
DECLARE sortOrder INTEGER;
BEGIN
    sortOrder := (SELECT coalesce(MAX(sort_order), 0) FROM storyboard_column WHERE goal_id = goalId) + 1;
    INSERT INTO storyboard_column (storyboard_id, goal_id, sort_order) VALUES (storyBoardId, goalId, sortOrder);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyBoardId;
END;
$$;

-- Revise a Storyboard Column --
CREATE PROCEDURE revise_storyboard_column(storyBoardId UUID, columnId UUID, columnName VARCHAR(256))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard_column SET name = columnName, updated_date = NOW() WHERE id = columnId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyBoardId;
END;
$$;

-- Delete a Storyboard Column --
CREATE PROCEDURE delete_storyboard_column(columnId UUID)
LANGUAGE plpgsql AS $$
DECLARE goalId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT goal_id, sort_order INTO goalId, sortOrder FROM storyboard_column WHERE id = columnId;

    DELETE FROM storyboard_story WHERE column_id = columnId;
    DELETE FROM storyboard_column WHERE id = columnId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard_column sc SET sort_order = (sc.sort_order - 1) WHERE sc.goal_id = goalId AND sc.sort_order > sortOrder;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Create a Storyboard Story --
CREATE PROCEDURE create_storyboard_story(storyBoardId UUID, goalId UUID, columnId UUID)
LANGUAGE plpgsql AS $$
DECLARE sortOrder INTEGER;
BEGIN
    sortOrder := (SELECT coalesce(MAX(sort_order), 0) FROM storyboard_story WHERE columnId = columnId) + 1;
    INSERT INTO storyboard_story (storyboard_id, goal_id, column_id, sort_order) VALUES (storyBoardId, goalId, columnId, sortOrder);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyBoardId;
END;
$$;

-- Revise a Storyboard Story Name --
CREATE PROCEDURE update_story_name(storyId UUID, storyName VARCHAR(256))
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET name = storyName, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Revise a Storyboard Story Content --
CREATE PROCEDURE update_story_content(storyId UUID, storyContent TEXT)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET content = storyContent, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Revise a Storyboard Story Color --
CREATE PROCEDURE update_story_color(storyId UUID, storyColor VARCHAR(32))
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET color = storyColor, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Revise a Storyboard Story Points --
CREATE PROCEDURE update_story_points(storyId UUID, updatedPoints INTEGER)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET points = updatedPoints, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Revise a Storyboard Story Closed status --
CREATE PROCEDURE update_story_closed(storyId UUID, isClosed BOOL)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET closed = isClosed, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Move a Storyboard Story to a new column and/or goal --
CREATE PROCEDURE move_story(storyId UUID, goalId UUID, columnId UUID, placeBefore TEXT)
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

-- Delete a Storyboard Story --
CREATE PROCEDURE delete_storyboard_story(storyId UUID)
LANGUAGE plpgsql AS $$
DECLARE columnId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT column_id, sort_order, storyboard_id INTO columnId, sortOrder, storyboardId FROM storyboard_story WHERE id = storyId;
    DELETE FROM storyboard_story WHERE id = storyId;
    UPDATE storyboard_story ss SET sort_order = (ss.sort_order - 1) WHERE ss.column_id = columnId AND ss.sort_order > sortOrder;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Add a comment to Storyboard Story --
CREATE PROCEDURE story_comment_add(storyboardId UUID, storyId UUID, userId UUID, comment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO storyboard_story_comment (storyboard_id, story_id, user_id, comment) VALUES (storyboardId, storyId, userId, comment);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Add a Persona to Storyboard --
CREATE PROCEDURE persona_add(storyboardId UUID, personaName VARCHAR(256), personaRole VARCHAR(256), personaDescription TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO storyboard_persona (storyboard_id, name, role, description) VALUES (storyboardId, personaName, personaRole, personaDescription);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Edit a Storyboard Persona --
CREATE PROCEDURE persona_edit(storyboardId UUID, personaId UUID, personaName VARCHAR(256), personaRole VARCHAR(256), personaDescription TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE storyboard_persona SET name = personaName, role = personaRole, description = personaDescription, updated_date = NOW() WHERE id = personaId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Delete a Storyboard Persona --
CREATE PROCEDURE persona_delete(storyboardId UUID, personaId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_persona WHERE id = personaId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;

    COMMIT;
END;
$$;

-- Clean up Storyboards older than X Days --
CREATE PROCEDURE clean_storyboards(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;