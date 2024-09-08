-- +goose Up
-- +goose StatementBegin

-- Create an enum type for predefined estimation scales
CREATE TYPE thunderdome.estimation_scale_type AS ENUM (
    'fibonacci',
    'modified_fibonacci',
    't_shirt',
    'powers_of_two',
    'thunderdome_default',
    'custom'
);

-- Create a table for custom estimation scales
CREATE TABLE thunderdome.estimation_scale (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    scale_type thunderdome.estimation_scale_type NOT NULL,
    values TEXT[] NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    default_scale BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_default_scale UNIQUE (organization_id, team_id, default_scale),
    CONSTRAINT check_ownership_and_default CHECK (
        (is_public = true AND organization_id IS NULL AND team_id IS NULL) OR
        (is_public = false AND (organization_id IS NOT NULL OR team_id IS NULL)) OR
        (default_scale = true AND (organization_id IS NOT NULL OR team_id IS NOT NULL))
    )
);

-- Add indexes for better query performance
CREATE INDEX idx_estimation_scale_organization ON thunderdome.estimation_scale(organization_id);
CREATE INDEX idx_estimation_scale_team ON thunderdome.estimation_scale(team_id);
CREATE INDEX idx_estimation_scale_public ON thunderdome.estimation_scale(is_public);

-- Modify the thunderdome.poker table
ALTER TABLE thunderdome.poker
ADD COLUMN estimation_scale_id UUID REFERENCES thunderdome.estimation_scale(id);

-- Alter the points column to allow 8 character length instead of 3
ALTER TABLE thunderdome.poker_story ALTER COLUMN points TYPE VARCHAR(8);

-- Remove the default value from point_values_allowed
ALTER TABLE thunderdome.poker
ALTER COLUMN point_values_allowed DROP DEFAULT;

-- Create a function to convert JSONB to TEXT[]
CREATE OR REPLACE FUNCTION jsonb_to_text_array(j jsonb) RETURNS TEXT[] AS $$
DECLARE
    result TEXT[];
BEGIN
    SELECT array_agg(value::text)
    FROM jsonb_array_elements_text(j)
    INTO result;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

-- Convert JSONB array to TEXT array
ALTER TABLE thunderdome.poker
ALTER COLUMN point_values_allowed TYPE TEXT[]
USING jsonb_to_text_array(point_values_allowed);

-- Add a new default value for point_values_allowed
ALTER TABLE thunderdome.poker
ALTER COLUMN point_values_allowed SET DEFAULT ARRAY['1', '2', '3', '5', '8', '13', '?'];

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_estimation_scale_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically update the updated_at timestamp
CREATE TRIGGER update_estimation_scale_timestamp
BEFORE UPDATE ON thunderdome.estimation_scale
FOR EACH ROW
EXECUTE FUNCTION update_estimation_scale_timestamp();

-- Create a function to ensure only one default scale per organization/team
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

-- Create a trigger to enforce the single default scale rule
CREATE TRIGGER enforce_single_default_scale
BEFORE INSERT OR UPDATE ON thunderdome.estimation_scale
FOR EACH ROW
WHEN (NEW.default_scale = true)
EXECUTE FUNCTION ensure_single_default_scale();

-- Insert default estimation scales
INSERT INTO thunderdome.estimation_scale (name, description, scale_type, values, is_public, default_scale)
VALUES
('Thunderdome Default', 'Thunderdome default sequence for estimation', 'thunderdome_default'::thunderdome.estimation_scale_type, ARRAY['0', '1/2', '1', '2', '3', '5', '8', '13', '20', '21', '34', '55', '100', '?', '☕️'], true, true),
('Fibonacci', 'Standard Fibonacci sequence for estimation', 'fibonacci'::thunderdome.estimation_scale_type, ARRAY['0', '1', '2', '3', '5', '8', '13', '21', '34', '55', '89', '?'], true, false),
('Modified Fibonacci', 'Modified Fibonacci sequence for estimation', 'modified_fibonacci'::thunderdome.estimation_scale_type, ARRAY['0', '1', '2', '3', '5', '8', '13', '21', '40', '80', '100', '?'], true, false),
('T-Shirt Sizes', 'T-shirt size scale for rough estimation', 't_shirt'::thunderdome.estimation_scale_type, ARRAY['XXS','XS', 'S', 'M', 'L', 'XL', 'XXL', '?'], true, false),
('Powers of Two', 'Powers of two scale for exponential complexity', 'powers_of_two'::thunderdome.estimation_scale_type, ARRAY['1', '2', '4', '8', '16', '32', '64', '?'], true, false);

-- Update existing poker sessions to use the default Fibonacci scale
-- and
-- Update the point_values_allowed column to be an array of text values

UPDATE thunderdome.poker
SET estimation_scale_id = (SELECT id FROM thunderdome.estimation_scale WHERE name = 'Thunderdome Default')
WHERE estimation_scale_id IS NULL;

DROP FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, teamid uuid, OUT pokerid uuid);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the foreign key from the thunderdome.poker table
ALTER TABLE thunderdome.poker
DROP COLUMN IF EXISTS estimation_scale_id,
ALTER COLUMN point_values_allowed TYPE JSONB USING to_jsonb(point_values_allowed);

-- Remove the triggers
DROP TRIGGER IF EXISTS update_estimation_scale_timestamp ON thunderdome.estimation_scale;
DROP TRIGGER IF EXISTS enforce_single_default_scale ON thunderdome.estimation_scale;

-- Remove the functions
DROP FUNCTION IF EXISTS update_estimation_scale_timestamp();
DROP FUNCTION IF EXISTS ensure_single_default_scale();

-- Drop the estimation_scale table
DROP TABLE IF EXISTS thunderdome.estimation_scale;

CREATE OR REPLACE FUNCTION thunderdome.poker_create(leaderid uuid, pokername character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, hidevoteridentity boolean, joincode text, leadercode text, teamid uuid, OUT pokerid uuid)
 RETURNS uuid
 LANGUAGE plpgsql
AS $function$
BEGIN
    INSERT INTO thunderdome.poker (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, hide_voter_identity, join_code, leader_code, team_id)
        VALUES (leaderid, pokername, pointsAllowed, autoVoting, pointAverageRounding, hideVoterIdentity, joinCode, leaderCode, teamid)
        RETURNING id INTO pokerid;
    INSERT INTO thunderdome.poker_facilitator (poker_id, user_id) VALUES (pokerid, leaderid);
    INSERT INTO thunderdome.poker_user (poker_id, user_id) VALUES (pokerid, leaderid);
END;
$function$;

-- +goose StatementEnd
