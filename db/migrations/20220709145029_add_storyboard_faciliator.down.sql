DROP TABLE storyboard_facilitator;
ALTER TABLE retro_facilitator DROP CONSTRAINT fk_retro_facilitator_retro_id;
ALTER TABLE retro_facilitator DROP CONSTRAINT fk_retro_facilitator_user_id;

ALTER TABLE storyboard_story DROP COLUMN annotations;
ALTER TABLE storyboard_story DROP COLUMN link;
DROP TABLE storyboard_goal_persona;
DROP TABLE storyboard_column_persona;

CREATE OR REPLACE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code) VALUES (ownerId, storyboardName, joinCode) RETURNING id INTO storyId;

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE sb_facilitator_add(UUID, UUID);
DROP PROCEDURE sb_facilitator_remove(UUID, UUID);
DROP PROCEDURE sb_story_link_edit(UUID, TEXT);
DROP PROCEDURE sb_story_annotations_edit(UUID, JSONB);
DROP PROCEDURE sb_goal_persona_add(UUID, UUID, UUID);
DROP PROCEDURE sb_goal_persona_remove(UUID, UUID, UUID);
DROP PROCEDURE sb_column_persona_add(UUID, UUID, UUID);
DROP PROCEDURE sb_column_persona_remove(UUID, UUID, UUID);

-- Get a Storyboards Goals --
CREATE OR REPLACE FUNCTION get_storyboard_goals(storyboardId UUID) RETURNS table (
    id UUID, sort_order INTEGER, name VARCHAR(256), columns JSON
) AS $$
BEGIN
    RETURN QUERY
        SELECT
            sg.id,
            sg.sort_order,
            sg.name,
            COALESCE(json_agg(to_jsonb(t) - 'goal_id' ORDER BY t.sort_order) FILTER (WHERE t.id IS NOT NULL), '[]') AS columns
        FROM storyboard_goal sg
        LEFT JOIN (
            SELECT
                sc.*,
                COALESCE(
                    json_agg(stss ORDER BY stss.sort_order) FILTER (WHERE stss.id IS NOT NULL), '[]'
                ) AS stories
            FROM storyboard_column sc
            LEFT JOIN (
                SELECT
                    ss.*,
                    COALESCE(
                        json_agg(stcm ORDER BY stcm.created_date) FILTER (WHERE stcm.id IS NOT NULL), '[]'
                    ) AS comments
                FROM storyboard_story ss
                LEFT JOIN storyboard_story_comment stcm ON stcm.story_id = ss.id
                GROUP BY ss.id
            ) stss ON stss.column_id = sc.id
            GROUP BY sc.id
        ) t ON t.goal_id = sg.id
        WHERE sg.storyboard_id = storyboardId
        GROUP BY sg.id
        ORDER BY sg.sort_order;
END;
$$ LANGUAGE plpgsql;