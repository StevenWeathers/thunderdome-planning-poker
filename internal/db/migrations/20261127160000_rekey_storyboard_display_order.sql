-- +goose Up
-- Display order keys move from a base95 to a base62 alphabet (internal/fracdex), which rejects some existing keys.
-- Valid keys sort the same in both, so only lists holding an invalid key are renumbered a0, a1, ... in their
-- current order. The unique constraints are dropped meanwhile as a new key can match a not yet updated one.

-- +goose StatementBegin
-- Mirrors validateOrderKey in internal/fracdex
CREATE FUNCTION thunderdome.is_valid_display_order(key text) RETURNS boolean
    LANGUAGE plpgsql IMMUTABLE
AS $$
DECLARE
    head int;
    int_length int;
BEGIN
    IF key IS NULL OR key COLLATE "C" !~ '^[A-Za-z][0-9A-Za-z]*$' OR key = 'A' || repeat('0', 26) THEN
        RETURN false;
    END IF;

    -- the head character encodes the length of the integer part
    head := ascii(key);
    int_length := CASE WHEN head >= ascii('a') THEN head - ascii('a') + 2 ELSE ascii('Z') - head + 2 END;

    -- the fractional part must not end in zero
    RETURN length(key) = int_length OR (length(key) > int_length AND right(key, 1) <> '0');
END;
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Returns the nth (zero based) integer key, as fracdex generates when appending to an empty list
CREATE FUNCTION thunderdome.nth_display_order(n bigint) RETURNS text
    LANGUAGE plpgsql IMMUTABLE STRICT
AS $$
DECLARE
    digits CONSTANT text := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
    remaining bigint := n;
    capacity bigint := 62;
    width int := 1;
    result text := '';
BEGIN
    WHILE remaining >= capacity LOOP
        remaining := remaining - capacity;
        capacity := capacity * 62;
        width := width + 1;
    END LOOP;

    FOR i IN 1..width LOOP
        result := substr(digits, (remaining % 62)::int + 1, 1) || result;
        remaining := remaining / 62;
    END LOOP;

    RETURN chr(ascii('a') + width - 1) || result;
END;
$$;
-- +goose StatementEnd

ALTER TABLE thunderdome.storyboard_goal DROP CONSTRAINT storyboard_goal_storyboard_id_display_order_key;
UPDATE thunderdome.storyboard_goal t
SET display_order = thunderdome.nth_display_order(o.position)
FROM (
    SELECT id, row_number() OVER (PARTITION BY storyboard_id ORDER BY display_order, created_date, id) - 1 AS position
    FROM thunderdome.storyboard_goal
    WHERE storyboard_id IN (
        SELECT storyboard_id FROM thunderdome.storyboard_goal
        WHERE NOT thunderdome.is_valid_display_order(display_order)
    )
) o
WHERE t.id = o.id;
ALTER TABLE thunderdome.storyboard_goal
    ADD CONSTRAINT storyboard_goal_storyboard_id_display_order_key UNIQUE (storyboard_id, display_order);

ALTER TABLE thunderdome.storyboard_column DROP CONSTRAINT storyboard_column_goal_id_display_order_key;
UPDATE thunderdome.storyboard_column t
SET display_order = thunderdome.nth_display_order(o.position)
FROM (
    SELECT id, row_number() OVER (PARTITION BY goal_id ORDER BY display_order, created_date, id) - 1 AS position
    FROM thunderdome.storyboard_column
    WHERE goal_id IN (
        SELECT goal_id FROM thunderdome.storyboard_column
        WHERE NOT thunderdome.is_valid_display_order(display_order)
    )
) o
WHERE t.id = o.id;
ALTER TABLE thunderdome.storyboard_column
    ADD CONSTRAINT storyboard_column_goal_id_display_order_key UNIQUE (goal_id, display_order);

ALTER TABLE thunderdome.storyboard_story DROP CONSTRAINT storyboard_story_column_id_display_order_key;
UPDATE thunderdome.storyboard_story t
SET display_order = thunderdome.nth_display_order(o.position)
FROM (
    SELECT id, row_number() OVER (PARTITION BY column_id ORDER BY display_order, created_date, id) - 1 AS position
    FROM thunderdome.storyboard_story
    WHERE column_id IN (
        SELECT column_id FROM thunderdome.storyboard_story
        WHERE NOT thunderdome.is_valid_display_order(display_order)
    )
) o
WHERE t.id = o.id;
ALTER TABLE thunderdome.storyboard_story
    ADD CONSTRAINT storyboard_story_column_id_display_order_key UNIQUE (column_id, display_order);

DROP FUNCTION thunderdome.nth_display_order(bigint);
DROP FUNCTION thunderdome.is_valid_display_order(text);

-- +goose Down
-- The new keys keep the same order and are also valid for the previous library, so there is nothing to restore.
SELECT 1;
