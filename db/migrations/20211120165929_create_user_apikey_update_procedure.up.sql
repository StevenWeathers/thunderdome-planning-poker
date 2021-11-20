-- Updates a user apikey
CREATE OR REPLACE PROCEDURE user_apikey_update(
    apikeyId text,
    userId uuid,
    keyActive boolean
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE api_keys SET active = keyActive, updated_date = NOW() WHERE id = apikeyId AND user_id = userId;
    UPDATE users SET last_active = NOW() WHERE id = userId;

    COMMIT;
END;
$$;