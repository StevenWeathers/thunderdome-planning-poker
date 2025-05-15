-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.user_email_change (
    change_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '01:00:00'::interval)
);

-- Create the function for the trigger
CREATE OR REPLACE FUNCTION thunderdome.delete_expired_user_email_changes()
RETURNS TRIGGER AS $$
DECLARE
    row_count INT;
BEGIN
    -- Delete expired records
    DELETE FROM thunderdome.user_email_change
    WHERE expire_date < NOW();

    -- Get the number of affected rows
    GET DIAGNOSTICS row_count = ROW_COUNT;

    -- Log only if records were deleted
    IF row_count > 0 THEN
        RAISE INFO 'Deleted % expired row(s) from thunderdome.user_email_change', row_count;
    END IF;

    RETURN NULL; -- for AFTER triggers
END;
$$ LANGUAGE plpgsql;

-- Create the trigger (only for INSERT and UPDATE)
CREATE TRIGGER trigger_delete_expired_user_email_changes
AFTER INSERT OR UPDATE ON thunderdome.user_email_change
FOR EACH STATEMENT
EXECUTE FUNCTION thunderdome.delete_expired_user_email_changes();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the trigger
DROP TRIGGER IF EXISTS trigger_delete_expired_user_email_changes ON thunderdome.user_email_change;
-- Drop the function
DROP FUNCTION IF EXISTS thunderdome.delete_expired_user_email_changes();
-- Drop the table
DROP TABLE IF EXISTS thunderdome.user_email_change;
-- +goose StatementEnd
