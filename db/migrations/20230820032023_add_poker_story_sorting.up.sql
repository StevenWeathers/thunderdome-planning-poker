ALTER TABLE thunderdome.poker_story ADD COLUMN position DOUBLE PRECISION DEFAULT 0 NOT NULL;

DO $$ DECLARE
    story RECORD;
    pid UUID;
    pos DOUBLE PRECISION = 0;
BEGIN
    FOR story IN SELECT id, poker_id, created_date FROM thunderdome.poker_story ORDER BY poker_id, created_date
    LOOP
        IF pid IS NULL OR story.poker_id <> pid THEN
            pid = story.poker_id;
            pos = -1;
        END IF;
        pos = pos + 1;

        UPDATE thunderdome.poker_story SET position = pos WHERE id = story.id;
    END LOOP;
END$$;

ALTER TABLE thunderdome.poker_story ADD CONSTRAINT poker_story_poker_id_position UNIQUE (poker_id, position);

CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
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
END$procedure$;