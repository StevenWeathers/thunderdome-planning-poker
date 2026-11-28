-- +goose Up
-- The previous triggers ran once per statement, where NEW and OLD are null, so they updated nothing.
-- Deletes are left out because removing a user cascades to these tables and must not revive old games.
-- last_active is only rewritten when over an hour old, which is precise enough for cleanup by age and
-- avoids locking the game row on every write, which could deadlock concurrent updates to a game.
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_poker_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE thunderdome.poker SET last_active = NOW()
    WHERE id = NEW.poker_id AND last_active < NOW() - interval '1 hour';
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_retro_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE thunderdome.retro SET last_active = NOW()
    WHERE id = NEW.retro_id AND last_active < NOW() - interval '1 hour';
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_storyboard_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    UPDATE thunderdome.storyboard SET last_active = NOW()
    WHERE id = NEW.storyboard_id AND last_active < NOW() - interval '1 hour';
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS poker_facilitator_poker_last_active ON thunderdome.poker_facilitator;
CREATE TRIGGER poker_facilitator_poker_last_active AFTER INSERT OR UPDATE ON thunderdome.poker_facilitator
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS poker_story_poker_last_active ON thunderdome.poker_story;
CREATE TRIGGER poker_story_poker_last_active AFTER INSERT OR UPDATE ON thunderdome.poker_story
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS poker_user_poker_last_active ON thunderdome.poker_user;
CREATE TRIGGER poker_user_poker_last_active AFTER INSERT OR UPDATE ON thunderdome.poker_user
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS retro_action_retro_last_active ON thunderdome.retro_action;
CREATE TRIGGER retro_action_retro_last_active AFTER INSERT OR UPDATE ON thunderdome.retro_action
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_facilitator_retro_last_active ON thunderdome.retro_facilitator;
CREATE TRIGGER retro_facilitator_retro_last_active AFTER INSERT OR UPDATE ON thunderdome.retro_facilitator
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_group_retro_last_active ON thunderdome.retro_group;
CREATE TRIGGER retro_group_retro_last_active AFTER INSERT OR UPDATE ON thunderdome.retro_group
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_item_retro_last_active ON thunderdome.retro_item;
CREATE TRIGGER retro_item_retro_last_active AFTER INSERT OR UPDATE ON thunderdome.retro_item
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_user_retro_last_active ON thunderdome.retro_user;
CREATE TRIGGER retro_user_retro_last_active AFTER INSERT OR UPDATE ON thunderdome.retro_user
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS storyboard_facilitator_storyboard_last_active ON thunderdome.storyboard_facilitator;
CREATE TRIGGER storyboard_facilitator_storyboard_last_active AFTER INSERT OR UPDATE ON thunderdome.storyboard_facilitator
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
DROP TRIGGER IF EXISTS storyboard_story_storyboard_last_active ON thunderdome.storyboard_story;
CREATE TRIGGER storyboard_story_storyboard_last_active AFTER INSERT OR UPDATE ON thunderdome.storyboard_story
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
DROP TRIGGER IF EXISTS storyboard_user_storyboard_last_active ON thunderdome.storyboard_user;
CREATE TRIGGER storyboard_user_storyboard_last_active AFTER INSERT OR UPDATE ON thunderdome.storyboard_user
    FOR EACH ROW EXECUTE FUNCTION thunderdome.update_storyboard_last_active();

-- +goose Down
DROP TRIGGER IF EXISTS poker_facilitator_poker_last_active ON thunderdome.poker_facilitator;
CREATE TRIGGER poker_facilitator_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_facilitator
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS poker_story_poker_last_active ON thunderdome.poker_story;
CREATE TRIGGER poker_story_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_story
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS poker_user_poker_last_active ON thunderdome.poker_user;
CREATE TRIGGER poker_user_poker_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.poker_user
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_poker_last_active();
DROP TRIGGER IF EXISTS retro_action_retro_last_active ON thunderdome.retro_action;
CREATE TRIGGER retro_action_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_action
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_facilitator_retro_last_active ON thunderdome.retro_facilitator;
CREATE TRIGGER retro_facilitator_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_facilitator
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_group_retro_last_active ON thunderdome.retro_group;
CREATE TRIGGER retro_group_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_group
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_item_retro_last_active ON thunderdome.retro_item;
CREATE TRIGGER retro_item_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_item
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS retro_user_retro_last_active ON thunderdome.retro_user;
CREATE TRIGGER retro_user_retro_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.retro_user
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_retro_last_active();
DROP TRIGGER IF EXISTS storyboard_facilitator_storyboard_last_active ON thunderdome.storyboard_facilitator;
CREATE TRIGGER storyboard_facilitator_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_facilitator
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
DROP TRIGGER IF EXISTS storyboard_story_storyboard_last_active ON thunderdome.storyboard_story;
CREATE TRIGGER storyboard_story_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_story
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();
DROP TRIGGER IF EXISTS storyboard_user_storyboard_last_active ON thunderdome.storyboard_user;
CREATE TRIGGER storyboard_user_storyboard_last_active AFTER INSERT OR DELETE OR UPDATE ON thunderdome.storyboard_user
    FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.update_storyboard_last_active();

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_poker_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.poker SET last_active = NOW() WHERE id = OLD.poker_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.poker SET last_active = NOW() WHERE id = NEW.poker_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.poker SET last_active = NOW() WHERE id = NEW.poker_id;
        RETURN NEW;
    END IF;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_retro_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.retro SET last_active = NOW() WHERE id = OLD.retro_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.retro SET last_active = NOW() WHERE id = NEW.retro_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.retro SET last_active = NOW() WHERE id = NEW.retro_id;
        RETURN NEW;
    END IF;
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION thunderdome.update_storyboard_last_active() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
       UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = OLD.storyboard_id;
       RETURN OLD;
    ELSIF TG_OP = 'UPDATE' THEN
       UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = NEW.storyboard_id;

       RETURN NEW;
    ELSIF TG_OP = 'INSERT' THEN
        UPDATE thunderdome.storyboard SET last_active = NOW() WHERE id = NEW.storyboard_id;
        RETURN NEW;
    END IF;
END;
$$;
-- +goose StatementEnd
