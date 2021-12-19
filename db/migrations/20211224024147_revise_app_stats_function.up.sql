-- Get Application Stats e.g. total user and battle counts
DROP FUNCTION get_app_stats();
CREATE FUNCTION get_app_stats(
    OUT unregistered_user_count INTEGER,
    OUT registered_user_count INTEGER,
    OUT battle_count INTEGER,
    OUT plan_count INTEGER,
    OUT organization_count INTEGER,
    OUT department_count INTEGER,
    OUT team_count INTEGER,
    OUT apikey_count INTEGER,
    OUT active_battle_count INTEGER,
    OUT active_battle_user_count INTEGER,
    OUT team_checkins_count INTEGER,
    OUT retro_count INTEGER,
    OUT active_retro_count INTEGER,
    OUT active_retro_user_count INTEGER,
    OUT retro_item_count INTEGER,
    OUT retro_action_count INTEGER
) AS $$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO battle_count FROM battles;
    SELECT COUNT(*) INTO plan_count FROM plans;
    SELECT COUNT(*) INTO organization_count FROM organization;
    SELECT COUNT(*) INTO department_count FROM organization_department;
    SELECT COUNT(*) INTO team_count FROM team;
    SELECT COUNT(*) INTO apikey_count FROM api_keys;
    SELECT COUNT(DISTINCT battle_id), COUNT(user_id)
        INTO active_battle_count, active_battle_user_count
        FROM battles_users WHERE active IS true;
    SELECT COUNT(*) INTO team_checkins_count FROM team_checkin;
    SELECT COUNT(*) INTO retro_count FROM retro;
    SELECT COUNT(DISTINCT retro_id), COUNT(user_id)
        INTO active_retro_count, active_retro_user_count
        FROM retro_user WHERE active IS true;
    SELECT COUNT(*) INTO retro_item_count FROM retro_item;
    SELECT COUNT(*) INTO retro_action_count FROM retro_action;
END;
$$ LANGUAGE plpgsql;