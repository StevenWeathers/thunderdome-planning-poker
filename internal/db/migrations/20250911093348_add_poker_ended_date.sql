-- +goose Up
-- +goose StatementBegin

-- Add ended_date column to poker table to track when games are stopped/ended
ALTER TABLE thunderdome.poker
ADD COLUMN ended_date TIMESTAMP WITH TIME ZONE;

-- Create index for performance when querying ended games
CREATE INDEX idx_poker_ended_date ON thunderdome.poker(ended_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the index first
DROP INDEX IF EXISTS thunderdome.idx_poker_ended_date;

-- Remove the ended_date column
ALTER TABLE thunderdome.poker
DROP COLUMN IF EXISTS ended_date;

-- +goose StatementEnd