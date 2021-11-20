-- Insert a new user apikey
CREATE OR REPLACE FUNCTION user_apikey_add(
    IN apikeyId text,
    IN keyName VARCHAR(256),
    IN userId uuid,
    OUT createdDate timestamp
)
AS $$
BEGIN
    INSERT INTO api_keys (id, name, user_id) VALUES (apikeyId, keyName, userId) RETURNING created_date INTO createdDate;
    UPDATE users SET last_active = NOW() WHERE id = userId;
END;
$$ LANGUAGE plpgsql;