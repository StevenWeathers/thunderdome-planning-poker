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
            pos = 0;
        END IF;
        pos = pos + 1;

        UPDATE thunderdome.poker_story SET position = pos WHERE id = story.id;
    END LOOP;
END$$;

ALTER TABLE thunderdome.poker_story ADD CONSTRAINT poker_story_poker_id_position UNIQUE (poker_id, position);