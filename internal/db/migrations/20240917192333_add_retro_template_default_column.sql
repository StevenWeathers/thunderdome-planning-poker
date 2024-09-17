-- +goose Up
-- +goose StatementBegin

-- Add the default_template column
ALTER TABLE thunderdome.retro_template
ADD COLUMN default_template BOOLEAN NOT NULL DEFAULT FALSE;

-- Create a function to ensure only one default template per organization, team, or public
CREATE OR REPLACE FUNCTION enforce_single_default_template()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.default_template = TRUE THEN
        -- If it's a public template
        IF NEW.is_public = TRUE THEN
            UPDATE thunderdome.retro_template
            SET default_template = FALSE
            WHERE is_public = TRUE AND id != NEW.id;
        -- If it's an organization template
        ELSIF NEW.organization_id IS NOT NULL THEN
            UPDATE thunderdome.retro_template
            SET default_template = FALSE
            WHERE organization_id = NEW.organization_id AND id != NEW.id;
        -- If it's a team template
        ELSIF NEW.team_id IS NOT NULL THEN
            UPDATE thunderdome.retro_template
            SET default_template = FALSE
            WHERE team_id = NEW.team_id AND id != NEW.id;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to enforce the single default template rule
CREATE TRIGGER enforce_single_default_template_trigger
BEFORE INSERT OR UPDATE ON thunderdome.retro_template
FOR EACH ROW
EXECUTE FUNCTION enforce_single_default_template();

-- Set a default template for public templates
UPDATE thunderdome.retro_template
SET default_template = TRUE
WHERE is_public = TRUE
AND id = (SELECT id FROM thunderdome.retro_template WHERE is_public = TRUE ORDER BY created_at LIMIT 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the trigger
DROP TRIGGER IF EXISTS enforce_single_default_template_trigger ON thunderdome.retro_template;

-- Remove the function
DROP FUNCTION IF EXISTS enforce_single_default_template();

-- Remove the default_template column
ALTER TABLE thunderdome.retro_template
DROP COLUMN IF EXISTS default_template;

-- +goose StatementEnd