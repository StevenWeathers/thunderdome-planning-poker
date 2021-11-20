DROP FUNCTION user_apikey_add(
    IN apikeyId text,
    IN keyName VARCHAR(256),
    IN userId uuid,
    OUT createdDate timestamp
);