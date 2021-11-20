DROP FUNCTION apikeys_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
);

DROP FUNCTION insert_user_reset(
    IN userEmail VARCHAR(320),
    OUT resetId UUID,
    OUT userId UUID,
    OUT userName VARCHAR(64)
);

DROP FUNCTION register_user(
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
);

DROP FUNCTION register_existing_user(
    IN activeUserId UUID,
    IN userName VARCHAR(64),
    IN userEmail VARCHAR(320),
    IN hashedPassword TEXT,
    IN userType VARCHAR(128),
    OUT userId UUID,
    OUT verifyId UUID
);

DROP FUNCTION registered_users_list(
    IN l_limit INTEGER,
    IN l_offset INTEGER
);

DROP FUNCTION registered_users_email_search(
    IN email_search VARCHAR(320),
    IN l_limit INTEGER,
    IN l_offset INTEGER
);