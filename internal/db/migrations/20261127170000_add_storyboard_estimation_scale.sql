-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.storyboard
ADD COLUMN estimation_scale_id UUID REFERENCES thunderdome.estimation_scale(id) ON DELETE SET NULL;

CREATE INDEX idx_storyboard_estimation_scale ON thunderdome.storyboard(estimation_scale_id);

-- match thunderdome.poker_story so that any scale value fits
ALTER TABLE thunderdome.storyboard_story ALTER COLUMN points TYPE VARCHAR(8);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.storyboard_story ALTER COLUMN points TYPE VARCHAR(3) USING LEFT(points, 3);

ALTER TABLE thunderdome.storyboard DROP COLUMN estimation_scale_id;
-- +goose StatementEnd
