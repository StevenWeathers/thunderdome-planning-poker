-- +goose Up
-- Moving a column between goals did not update its stories, leaving them tied to the old goal,
-- so deleting that goal would also delete them.
UPDATE thunderdome.storyboard_story s
SET goal_id = c.goal_id
FROM thunderdome.storyboard_column c
WHERE s.column_id = c.id AND s.goal_id IS DISTINCT FROM c.goal_id;

-- +goose Down
-- The previous goal_id values were wrong, so there is nothing to restore.
SELECT 1;
