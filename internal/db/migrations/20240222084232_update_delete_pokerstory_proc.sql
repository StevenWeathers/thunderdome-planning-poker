-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE
	active_storyid UUID;
BEGIN
    active_storyid := (SELECT b.active_story_id FROM thunderdome.poker b WHERE b.id = pokerid);
    DELETE FROM thunderdome.poker_story WHERE id = storyid;

    IF active_storyid = storyid THEN
        UPDATE thunderdome.poker SET last_active = NOW(), voting_locked = true, active_story_id = null
        WHERE id = pokerid;
    END IF;

    COMMIT;
END$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE
	active_storyid UUID;
	story RECORD;
    pos DOUBLE PRECISION = -1;
BEGIN
    active_storyid := (SELECT b.active_story_id FROM thunderdome.poker b WHERE b.id = pokerid);
    DELETE FROM thunderdome.poker_story WHERE id = storyid;

	FOR story IN SELECT id, position FROM thunderdome.poker_story WHERE poker_id = pokerid ORDER BY position
    LOOP
        pos = pos + 1;

        UPDATE thunderdome.poker_story SET position = pos WHERE id = story.id;
    END LOOP;

    IF active_storyid = storyid THEN
        UPDATE thunderdome.poker SET last_active = NOW(), voting_locked = true, active_story_id = null
        WHERE id = pokerid;
    END IF;

    COMMIT;
END$$;
-- +goose StatementEnd
