-- +goose Up
-- +goose StatementBegin

-- Drop the existing stored procedure
DROP PROCEDURE IF EXISTS thunderdome.user_password_reset(uuid, text);

-- Create the function for the trigger
CREATE OR REPLACE FUNCTION thunderdome.delete_expired_user_resets()
RETURNS TRIGGER AS $$
DECLARE
    row_count INT;
BEGIN
    -- Delete expired records
    DELETE FROM thunderdome.user_reset
    WHERE expire_date < NOW();

    -- Get the number of affected rows
    GET DIAGNOSTICS row_count = ROW_COUNT;

    -- Log only if records were deleted
    IF row_count > 0 THEN
        RAISE INFO 'Deleted % expired row(s) from thunderdome.user_reset', row_count;
    END IF;

    RETURN NULL; -- for AFTER triggers
END;
$$ LANGUAGE plpgsql;

-- Create the trigger (only for INSERT and UPDATE)
CREATE TRIGGER trigger_delete_expired_user_resets
AFTER INSERT OR UPDATE ON thunderdome.user_reset
FOR EACH STATEMENT
EXECUTE FUNCTION thunderdome.delete_expired_user_resets();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop the trigger
DROP TRIGGER IF EXISTS trigger_delete_expired_user_resets ON thunderdome.user_reset;

-- Drop the function
DROP FUNCTION IF EXISTS thunderdome.delete_expired_user_resets();

-- Recreate the original stored procedure
CREATE OR REPLACE PROCEDURE thunderdome.user_password_reset(IN resetid uuid, IN userpassword text)
LANGUAGE plpgsql
AS $procedure$
DECLARE matchedUserId UUID;
BEGIN
    matchedUserId := (
        SELECT w.id
        FROM thunderdome.user_reset wr
        LEFT JOIN thunderdome.users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );
    IF matchedUserId IS NULL THEN
        -- attempt delete in case reset record expired
        DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;
    UPDATE thunderdome.auth_credential SET password = userPassword, updated_date = NOW() WHERE user_id = matchedUserId;
    UPDATE thunderdome.users SET last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM thunderdome.user_reset WHERE reset_id = resetId;
    COMMIT;
END;
$procedure$;

-- +goose StatementEnd