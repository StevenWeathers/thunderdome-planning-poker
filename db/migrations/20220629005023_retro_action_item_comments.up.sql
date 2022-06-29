CREATE TABLE retro_action_comment (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    action_id UUID REFERENCES retro_action(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    comment TEXT,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE retro_action_assignee (
    action_id UUID REFERENCES retro_action(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (action_id, user_id)
);

-- Add a comment to Retro Action --
CREATE PROCEDURE retro_action_comment_add(retroId UUID, actionId UUID, userId UUID, actionComment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO retro_action_comment (action_id, user_id, comment) VALUES (actionId, userId, actionComment);
    UPDATE retro_action SET updated_date = NOW() WHERE id = actionId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;

    COMMIT;
END;
$$;

-- Edit a Retro Action comment --
CREATE PROCEDURE retro_action_comment_edit(retroId UUID, actionId UUID, commentId UUID, actionComment TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro_action_comment SET comment = actionComment WHERE id = commentId;
    UPDATE retro_action SET updated_date = NOW() WHERE id = actionId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;

    COMMIT;
END;
$$;

-- Delete a Retro Action comment --
CREATE PROCEDURE retro_action_comment_delete(retroId UUID, actionId UUID, commentId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retro_action_comment WHERE id = commentId;
    UPDATE retro_action SET updated_date = NOW() WHERE id = actionId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;

    COMMIT;
END;
$$;

-- Add an assignee to Retro Action --
CREATE PROCEDURE retro_action_assignee_add(retroId UUID, actionId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO retro_action_assignee (action_id, user_id) VALUES (actionId, userId);
    UPDATE retro_action SET updated_date = NOW() WHERE id = actionId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;

    COMMIT;
END;
$$;

-- Delete a Retro Action assignee --
CREATE PROCEDURE retro_action_assignee_delete(retroId UUID, actionId UUID, userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM retro_action_assignee WHERE action_id = actionId AND user_id = userId;
    UPDATE retro_action SET updated_date = NOW() WHERE id = actionId;
    UPDATE retro SET updated_date = NOW() WHERE id = retroId;

    COMMIT;
END;
$$;