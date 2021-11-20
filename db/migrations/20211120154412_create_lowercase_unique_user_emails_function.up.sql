-- Find and update all user emails that include an uppercase character and dont have a duplicate account to lowercase email --
CREATE OR REPLACE FUNCTION lowercase_unique_user_emails() RETURNS table (
    name VARCHAR(256), email VARCHAR(320)
) AS $$
BEGIN
    RETURN QUERY
        UPDATE users u
        SET email = lower(u.email), updated_date = NOW()
        FROM (
            SELECT lower(su.email) AS email
            FROM users su
            WHERE su.email IS NOT NULL
            GROUP BY lower(su.email) HAVING count(su.*) = 1
        ) AS subquery
        WHERE lower(u.email) = subquery.email AND u.email ~ '[A-Z]' RETURNING u.name, u.email;
END;
$$ LANGUAGE plpgsql;