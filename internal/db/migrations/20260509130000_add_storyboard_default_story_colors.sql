-- +goose Up
ALTER TABLE thunderdome.storyboard_goal ADD COLUMN IF NOT EXISTS default_story_color character varying(32);
ALTER TABLE thunderdome.storyboard_column ADD COLUMN IF NOT EXISTS default_story_color character varying(32);

-- +goose Down
ALTER TABLE thunderdome.storyboard_column DROP COLUMN IF EXISTS default_story_color;
ALTER TABLE thunderdome.storyboard_goal DROP COLUMN IF EXISTS default_story_color;