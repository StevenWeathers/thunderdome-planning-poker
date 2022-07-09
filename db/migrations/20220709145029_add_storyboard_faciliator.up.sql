ALTER TABLE retro_facilitator ADD CONSTRAINT fk_retro_facilitator_retro_id FOREIGN KEY (retro_id) REFERENCES retro(id) ON DELETE CASCADE;
ALTER TABLE retro_facilitator ADD CONSTRAINT fk_retro_facilitator_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE storyboard_story ADD COLUMN annotations jsonb NOT NULL DEFAULT '[]'::jsonb;
ALTER TABLE storyboard_story ADD COLUMN link TEXT;

CREATE TABLE "storyboard_goal_persona" (
    goal_id UUID REFERENCES storyboard_goal(id) ON DELETE CASCADE,
    persona_id UUID REFERENCES storyboard_persona(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ("goal_id", "persona_id")
);

CREATE TABLE "storyboard_column_persona" (
    column_id UUID REFERENCES storyboard_column(id) ON DELETE CASCADE,
    persona_id UUID REFERENCES storyboard_persona(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ("column_id", "persona_id")
);

CREATE TABLE "storyboard_facilitator" (
    "storyboard_id" UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    "user_id" UUID REFERENCES users(id) ON DELETE CASCADE,
    "created_date" TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ("storyboard_id","user_id")
);

-- back fill the storyboard_facilitator table with owner ids
BEGIN;
    INSERT INTO storyboard_facilitator (storyboard_id, user_id)
    SELECT id, owner_id FROM storyboard;
COMMIT;

CREATE OR REPLACE FUNCTION create_storyboard(ownerId UUID, storyboardName VARCHAR(256), joinCode VARCHAR(128)) RETURNS UUID
AS $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code) VALUES (ownerId, storyboardName, joinCode) RETURNING id INTO storyId;
    INSERT INTO storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);

    RETURN storyId;
END;
$$ LANGUAGE plpgsql;

-- add storyboard facilitator --
CREATE PROCEDURE sb_facilitator_add(storyboardId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO storyboard_facilitator (storyboard_id, user_id) VALUES (storyboardId, userId);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- remove storyboard facilitator --
CREATE PROCEDURE sb_facilitator_remove(storyboardId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_facilitator WHERE storyboard_id = storyboardId AND user_id = userId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- edit storyboard story link --
CREATE PROCEDURE sb_story_link_edit(storyId UUID, updatedLink TEXT)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET link = updatedLink, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- edit storyboard story annotations --
CREATE PROCEDURE sb_story_annotations_edit(storyId UUID, updatedAnnotations JSONB)
LANGUAGE plpgsql AS $$
DECLARE storyboardId UUID;
BEGIN
    UPDATE storyboard_story SET annotations = updatedAnnotations, updated_date = NOW() WHERE id = storyId RETURNING storyboard_id INTO storyboardId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- add storyboard goal persona --
CREATE PROCEDURE sb_goal_persona_add(storyboardId UUID, goalId UUID, personaId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO storyboard_goal_persona (goal_id, persona_id) VALUES (goalId, personaId);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- remove storyboard goal persona --
CREATE PROCEDURE sb_goal_persona_remove(storyboardId UUID, goalId UUID, personaId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_goal_persona WHERE goal_id = goalId AND persona_id = personaId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- add storyboard column persona --
CREATE PROCEDURE sb_column_persona_add(storyboardId UUID, columnId UUID, personaId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO storyboard_column_persona (column_id, persona_id) VALUES (columnId, personaId);
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- remove storyboard column persona --
CREATE PROCEDURE sb_column_persona_remove(storyboardId UUID, columnId UUID, personaId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM storyboard_column_persona WHERE column_id = columnId AND persona_id = personaId;
    UPDATE storyboard SET updated_date = NOW() WHERE id = storyboardId;
END;
$$;

-- Get a Storyboards Goals --
DROP FUNCTION get_storyboard_goals(UUID);
CREATE FUNCTION get_storyboard_goals(storyboardId UUID) RETURNS table (
    id UUID, sort_order INTEGER, name VARCHAR(256), columns JSON, personas JSON
) AS $$
BEGIN
    RETURN QUERY
        SELECT
            sg.id,
            sg.sort_order,
            sg.name,
            COALESCE(json_agg(to_jsonb(t) - 'goal_id' ORDER BY t.sort_order) FILTER (WHERE t.id IS NOT NULL), '[]') AS columns,
            COALESCE(json_agg(to_jsonb(sgp) - 'goal_id') FILTER (WHERE sgp.goal_id IS NOT NULL), '[]') AS personas
        FROM storyboard_goal sg
        LEFT JOIN (
            SELECT
                sc.*,
                COALESCE(
                    json_agg(stss ORDER BY stss.sort_order) FILTER (WHERE stss.id IS NOT NULL), '[]'
                ) AS stories,
                COALESCE(
                    json_agg(scp) FILTER (WHERE scp.column_id IS NOT NULL), '[]'
                ) AS personas
            FROM storyboard_column sc
            LEFT JOIN (
                SELECT cp.column_id, sp.*
                FROM storyboard_column_persona cp
                LEFT JOIN storyboard_persona sp ON sp.id = cp.persona_id
            ) scp ON scp.column_id = sc.id
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
        LEFT JOIN (
            SELECT gp.goal_id, sp.*
            FROM storyboard_goal_persona gp
            LEFT JOIN storyboard_persona sp ON sp.id = gp.persona_id
        ) sgp ON sgp.goal_id = sg.id
        WHERE sg.storyboard_id = storyboardId
        GROUP BY sg.id
        ORDER BY sg.sort_order;
END;
$$ LANGUAGE plpgsql;