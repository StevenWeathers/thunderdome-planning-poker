DROP PROCEDURE set_retro_owner(retroId UUID, ownerId UUID);
DROP PROCEDURE set_retro_phase(retroId UUID, nextPhase SMALLINT);
DROP PROCEDURE delete_retro(retroId UUID);
DROP FUNCTION create_retro(ownerId UUID, retroName VARCHAR(256), format VARCHAR(32), joinCode VARCHAR(128));
DROP FUNCTION get_retros_by_user(userId UUID);
DROP FUNCTION get_retro_users(retroId UUID);
DROP FUNCTION get_retro_user(retroId UUID, userId UUID);
DROP FUNCTION team_retro_list(IN teamId UUID, IN l_limit INTEGER, IN l_offset INTEGER);
DROP FUNCTION team_retro_add(IN teamId UUID, IN retroId UUID);
DROP FUNCTION team_retro_remove(IN teamId UUID, IN retroId UUID);
DROP PROCEDURE clean_retros(daysOld INTEGER);
CREATE OR REPLACE PROCEDURE deactivate_all_users()
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE battles_users SET active = false WHERE active = true;
END;
$$;