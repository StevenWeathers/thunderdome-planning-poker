-- +goose Up
-- +goose StatementBegin

-- Storyboard Story
ALTER TABLE thunderdome.storyboard_story ADD COLUMN display_order text COLLATE "C";
UPDATE thunderdome.storyboard_story SET display_order = 'a' || sort_order::text;
ALTER TABLE thunderdome.storyboard_story DROP CONSTRAINT storyboard_story_column_id_sort_order_key;
ALTER TABLE thunderdome.storyboard_story DROP COLUMN sort_order;
ALTER TABLE thunderdome.storyboard_story ADD CONSTRAINT storyboard_story_column_id_display_order_key UNIQUE (column_id, display_order);

-- Storyboard Column
ALTER TABLE thunderdome.storyboard_column ADD COLUMN display_order text COLLATE "C";
UPDATE thunderdome.storyboard_column SET display_order = 'a' || sort_order::text;
ALTER TABLE thunderdome.storyboard_column DROP CONSTRAINT storyboard_column_goal_id_sort_order_key;
ALTER TABLE thunderdome.storyboard_column DROP COLUMN sort_order;
ALTER TABLE thunderdome.storyboard_column ADD CONSTRAINT storyboard_column_goal_id_display_order_key UNIQUE (goal_id, display_order);

-- Storyboard Goal
ALTER TABLE thunderdome.storyboard_goal ADD COLUMN display_order text COLLATE "C";
UPDATE thunderdome.storyboard_goal SET display_order = 'a' || sort_order::text;
ALTER TABLE thunderdome.storyboard_goal DROP CONSTRAINT storyboard_goal_storyboard_id_sort_order_key;
ALTER TABLE thunderdome.storyboard_goal DROP COLUMN sort_order;
ALTER TABLE thunderdome.storyboard_goal ADD CONSTRAINT storyboard_goal_storyboard_id_display_order_key UNIQUE (storyboard_id, display_order);

DROP PROCEDURE "thunderdome".sb_column_delete(IN columnid uuid);
DROP PROCEDURE "thunderdome".sb_goal_delete(IN goalid uuid);
DROP PROCEDURE "thunderdome".sb_story_delete(IN storyid uuid);
DROP PROCEDURE "thunderdome".sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Storyboard Story
ALTER TABLE thunderdome.storyboard_story ADD COLUMN sort_order int4;
UPDATE thunderdome.storyboard_story SET sort_order = substring(display_order, 2)::int4;
ALTER TABLE thunderdome.storyboard_story DROP CONSTRAINT storyboard_story_column_id_display_order_key;
ALTER TABLE thunderdome.storyboard_story DROP COLUMN display_order;
ALTER TABLE thunderdome.storyboard_story ADD CONSTRAINT storyboard_story_column_id_sort_order_key UNIQUE (column_id, sort_order);

-- Storyboard Column
ALTER TABLE thunderdome.storyboard_column ADD COLUMN sort_order int4;
UPDATE thunderdome.storyboard_column SET sort_order = substring(display_order, 2)::int4;
ALTER TABLE thunderdome.storyboard_column DROP CONSTRAINT storyboard_column_goal_id_display_order_key;
ALTER TABLE thunderdome.storyboard_column DROP COLUMN display_order;
ALTER TABLE thunderdome.storyboard_column ADD CONSTRAINT storyboard_column_goal_id_sort_order_key UNIQUE (goal_id, sort_order);

-- Storyboard Goal
ALTER TABLE thunderdome.storyboard_goal ADD COLUMN sort_order int4;
UPDATE thunderdome.storyboard_goal SET sort_order = substring(display_order, 2)::int4;
ALTER TABLE thunderdome.storyboard_goal DROP CONSTRAINT storyboard_goal_storyboard_id_display_order_key;
ALTER TABLE thunderdome.storyboard_goal DROP COLUMN display_order;
ALTER TABLE thunderdome.storyboard_goal ADD CONSTRAINT storyboard_goal_storyboard_id_sort_order_key UNIQUE (storyboard_id, sort_order);

CREATE OR REPLACE PROCEDURE thunderdome.sb_story_move(IN storyid uuid, IN goalid uuid, IN columnid uuid, IN placebefore text)
 LANGUAGE plpgsql
AS $procedure$
DECLARE storyboardId UUID;
DECLARE srcGoalId UUID;
DECLARE srcColumnId UUID;
DECLARE srcSortOrder INTEGER;
DECLARE targetSortOrder INTEGER;
BEGIN
    SET CONSTRAINTS thunderdome.storyboard_story_column_id_sort_order_key DEFERRED;
    -- Get Story current details
    SELECT
        storyboard_id, goal_id, column_id, sort_order, name, color, content, created_date
    INTO
        storyboardId, srcGoalId, srcColumnId, srcSortOrder
    FROM thunderdome.storyboard_story WHERE id = storyId;

    -- Get target sort order
    IF placeBefore = '' THEN
        SELECT coalesce(max(sort_order), 0) + 1 INTO targetSortOrder FROM thunderdome.storyboard_story WHERE column_id = columnId;
    ELSE
        SELECT sort_order INTO targetSortOrder FROM thunderdome.storyboard_story WHERE column_id = columnId AND id = placeBefore::UUID;
    END IF;

    -- Remove from source column
    UPDATE thunderdome.storyboard_story SET column_id = columnId, sort_order = 9000 WHERE id = storyId;
    -- Update sort order in src column
    UPDATE thunderdome.storyboard_story ss SET sort_order = (t.sort_order - 1)
    FROM (
        SELECT id, sort_order FROM thunderdome.storyboard_story
        WHERE column_id = srcColumnId AND sort_order > srcSortOrder
        ORDER BY sort_order ASC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Update sort order for any story that should come after newly moved story
    UPDATE thunderdome.storyboard_story ss SET sort_order = (t.sort_order + 1)
    FROM (
        SELECT id, sort_order FROM thunderdome.storyboard_story
        WHERE column_id = columnId AND sort_order >= targetSortOrder
        ORDER BY sort_order DESC
        FOR UPDATE
    ) AS t
    WHERE ss.id = t.id;

    -- Finally, insert story in its ordered place
	UPDATE thunderdome.storyboard_story SET sort_order = targetSortOrder WHERE id = storyId;

    COMMIT;
END;
$procedure$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_story_delete(IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE columnId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT column_id, sort_order, storyboard_id INTO columnId, sortOrder, storyboardId
        FROM thunderdome.storyboard_story WHERE id = storyId;
    DELETE FROM thunderdome.storyboard_story WHERE id = storyId;
    UPDATE thunderdome.storyboard_story ss SET sort_order = (ss.sort_order - 1)
        WHERE ss.column_id = columnId AND ss.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_goal_delete(IN goalid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE storyboardId UUID;
DECLARE sortOrder INTEGER;
BEGIN
    SELECT sort_order, storyboard_id INTO sortOrder, storyboardId FROM thunderdome.storyboard_goal WHERE id = goalId;

    DELETE FROM thunderdome.storyboard_story WHERE goal_id = goalId;
    DELETE FROM thunderdome.storyboard_column WHERE goal_id = goalId;
    DELETE FROM thunderdome.storyboard_goal WHERE id = goalId;
    UPDATE thunderdome.storyboard_goal sg SET sort_order = (sg.sort_order - 1)
        WHERE sg.storyboard_id = storyBoardId AND sg.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

CREATE OR REPLACE PROCEDURE thunderdome.sb_column_delete(IN columnid uuid)
 LANGUAGE plpgsql
AS $procedure$
DECLARE goalId UUID;
DECLARE sortOrder INTEGER;
DECLARE storyboardId UUID;
BEGIN
    SELECT goal_id, sort_order INTO goalId, sortOrder FROM thunderdome.storyboard_column WHERE id = columnId;

    DELETE FROM thunderdome.storyboard_story WHERE column_id = columnId;
    DELETE FROM thunderdome.storyboard_column WHERE id = columnId RETURNING storyboard_id INTO storyboardId;
    UPDATE thunderdome.storyboard_column sc SET sort_order = (sc.sort_order - 1)
        WHERE sc.goal_id = goalId AND sc.sort_order > sortOrder;

    COMMIT;
END;
$procedure$;

-- +goose StatementEnd