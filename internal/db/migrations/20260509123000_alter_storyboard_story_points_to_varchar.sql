-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.storyboard_story
ALTER COLUMN points TYPE VARCHAR(3)
USING CASE
    WHEN points IS NULL THEN NULL
    ELSE points::VARCHAR(3)
END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.storyboard_story
ALTER COLUMN points TYPE INTEGER
USING CASE
    WHEN points IS NULL OR BTRIM(points) = '' THEN NULL
    WHEN BTRIM(points) ~ '^-?[0-9]+$' THEN BTRIM(points)::INTEGER
    ELSE NULL
END;
-- +goose StatementEnd