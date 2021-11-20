-- Deletes a user apikey
CREATE OR REPLACE PROCEDURE user_apikey_delete(
    apikeyId text,
    userId uuid
)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM api_keys WHERE id = apikeyId AND user_id = userId;
    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;