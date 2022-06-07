DROP FUNCTION IF EXISTS create_storyboard(UUID, VARCHAR);
DROP FUNCTION IF EXISTS get_storyboards_by_user(uuid);
DROP FUNCTION IF EXISTS get_storyboard_goals(uuid);
DROP FUNCTION IF EXISTS get_storyboard_users(uuid);
DROP FUNCTION IF EXISTS get_storyboard_personas(uuid);
DROP FUNCTION IF EXISTS get_storyboard_user(uuid, uuid);
DROP FUNCTION IF EXISTS team_storyboard_list(UUID,INTEGER,INTEGER);
DROP FUNCTION IF EXISTS team_storyboard_add(UUID,UUID);
DROP FUNCTION IF EXISTS team_storyboard_remove(UUID,UUID);