ALTER TABLE thunderdome.poker_story DROP COLUMN position;

CREATE OR REPLACE PROCEDURE thunderdome.poker_story_delete(IN pokerid uuid, IN storyid uuid)
 LANGUAGE plpgsql
AS $procedure$
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
END$procedure$;