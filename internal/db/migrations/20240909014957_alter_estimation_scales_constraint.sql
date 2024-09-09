-- +goose Up
-- +goose StatementBegin

-- Step 1: Drop the existing constraint
ALTER TABLE thunderdome.estimation_scale
DROP CONSTRAINT IF EXISTS check_ownership_and_default;

-- Step 2: Add back the original constraint
ALTER TABLE thunderdome.estimation_scale
ADD CONSTRAINT check_ownership_and_default CHECK (
    (is_public = true AND organization_id IS NULL AND team_id IS NULL) OR
    (is_public = false AND (organization_id IS NOT NULL OR team_id IS NOT NULL))
);

-- Step 3: Update the ensure_single_default_scale function
CREATE OR REPLACE FUNCTION ensure_single_default_scale()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.default_scale = true THEN
        -- For public scales, ensure no other public scale is set as default
        IF NEW.is_public = true THEN
            UPDATE thunderdome.estimation_scale
            SET default_scale = false
            WHERE is_public = true AND id != NEW.id;
        -- For organization-specific scales
        ELSIF NEW.organization_id IS NOT NULL THEN
            UPDATE thunderdome.estimation_scale
            SET default_scale = false
            WHERE organization_id = NEW.organization_id AND id != NEW.id;
        -- For team-specific scales
        ELSIF NEW.team_id IS NOT NULL THEN
            UPDATE thunderdome.estimation_scale
            SET default_scale = false
            WHERE team_id = NEW.team_id AND id != NEW.id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP FUNCTION thunderdome.appstats_get(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT poker_count integer, OUT poker_story_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_poker_count integer, OUT active_poker_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer);

DROP TYPE IF EXISTS thunderdome.usersvote;
CREATE TYPE thunderdome.usersvote AS (
    "warriorId" uuid,
    vote character varying(8)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Step 1: Drop the existing constraint
ALTER TABLE thunderdome.estimation_scale
DROP CONSTRAINT IF EXISTS check_ownership_and_default;

-- Step 2: Add back the original constraint
ALTER TABLE thunderdome.estimation_scale
ADD CONSTRAINT check_ownership_and_default CHECK (
    (is_public = true AND organization_id IS NULL AND team_id IS NULL) OR
    (is_public = false AND (organization_id IS NOT NULL OR team_id IS NULL)) OR
    (default_scale = true AND (organization_id IS NOT NULL OR team_id IS NOT NULL))
);

-- Step 3: Revert the ensure_single_default_scale function
CREATE OR REPLACE FUNCTION ensure_single_default_scale()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.default_scale = true THEN
        UPDATE thunderdome.estimation_scale
        SET default_scale = false
        WHERE (organization_id = NEW.organization_id OR team_id = NEW.team_id)
        AND id != NEW.id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION thunderdome.appstats_get(OUT unregistered_user_count integer, OUT registered_user_count integer, OUT poker_count integer, OUT poker_story_count integer, OUT organization_count integer, OUT department_count integer, OUT team_count integer, OUT apikey_count integer, OUT active_poker_count integer, OUT active_poker_user_count integer, OUT team_checkins_count integer, OUT retro_count integer, OUT active_retro_count integer, OUT active_retro_user_count integer, OUT retro_item_count integer, OUT retro_action_count integer, OUT storyboard_count integer, OUT active_storyboard_count integer, OUT active_storyboard_user_count integer, OUT storyboard_goal_count integer, OUT storyboard_column_count integer, OUT storyboard_story_count integer, OUT storyboard_persona_count integer)
 RETURNS record
 LANGUAGE plpgsql
AS $function$
BEGIN
    SELECT COUNT(*) INTO unregistered_user_count FROM thunderdome.users WHERE email IS NULL;
    SELECT COUNT(*) INTO registered_user_count FROM thunderdome.users WHERE email IS NOT NULL;
    SELECT COUNT(*) INTO poker_count FROM thunderdome.poker;
    SELECT COUNT(*) INTO poker_story_count FROM thunderdome.poker_story;
    SELECT COUNT(*) INTO organization_count FROM thunderdome.organization;
    SELECT COUNT(*) INTO department_count FROM thunderdome.organization_department;
    SELECT COUNT(*) INTO team_count FROM thunderdome.team;
    SELECT COUNT(*) INTO apikey_count FROM thunderdome.api_key;
    SELECT COUNT(DISTINCT poker_id), COUNT(user_id)
        INTO active_poker_count, active_poker_user_count
        FROM thunderdome.poker_user WHERE active IS true;
    SELECT COUNT(*) INTO team_checkins_count FROM thunderdome.team_checkin;
    SELECT COUNT(*) INTO retro_count FROM thunderdome.retro;
    SELECT COUNT(DISTINCT retro_id), COUNT(user_id)
        INTO active_retro_count, active_retro_user_count
        FROM thunderdome.retro_user WHERE active IS true;
    SELECT COUNT(*) INTO retro_item_count FROM thunderdome.retro_item;
    SELECT COUNT(*) INTO retro_action_count FROM thunderdome.retro_action;
    SELECT COUNT(*) INTO storyboard_count FROM thunderdome.storyboard;
    SELECT COUNT(DISTINCT storyboard_id), COUNT(user_id)
        INTO active_storyboard_count, active_storyboard_user_count
        FROM thunderdome.storyboard_user WHERE active IS true;
    SELECT COUNT(*) INTO storyboard_goal_count FROM thunderdome.storyboard_goal;
    SELECT COUNT(*) INTO storyboard_column_count FROM thunderdome.storyboard_column;
    SELECT COUNT(*) INTO storyboard_story_count FROM thunderdome.storyboard_story;
    SELECT COUNT(*) INTO storyboard_persona_count FROM thunderdome.storyboard_persona;
END;
$function$;

DROP TYPE IF EXISTS thunderdome.usersvote;
CREATE TYPE thunderdome.usersvote AS (
    "warriorId" uuid,
    vote character varying(3)
);
-- +goose StatementEnd
