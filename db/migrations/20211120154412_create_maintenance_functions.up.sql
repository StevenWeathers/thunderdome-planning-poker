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

-- Find and merge duplicate email user accounts caused by case sensitive bug --
CREATE OR REPLACE FUNCTION merge_nonunique_user_accounts() RETURNS table (
    name VARCHAR(256), email VARCHAR(320)
) AS $$
DECLARE usr RECORD;
BEGIN
    FOR usr IN
        SELECT
            array_agg(su.id) as id,
            array_agg(su.name) as name,
            lower(su.email) AS email,
            MAX(last_active) as active_date,
            array_agg(su.country) as country,
            array_agg(su.company) as company,
            array_agg(su.job_title) as job_title,
            array_agg(su.locale) as locale
        FROM users su
        WHERE su.email IS NOT NULL
        GROUP BY lower(su.email) HAVING count(su.*) > 1
        ORDER BY active_date DESC
    LOOP
        -- update battles
        UPDATE battles SET owner_id = usr.id[1] WHERE owner_id = usr.id[2];
        -- update battle_users
        BEGIN
            UPDATE battles_users SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in battle';
        END;
        -- update battle_leaders
        BEGIN
            UPDATE battles_leaders SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update organization_user
        BEGIN
            UPDATE organization_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update department_user
        BEGIN
            UPDATE department_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in department';
        END;
        -- update team_user
        BEGIN
            UPDATE team_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in team';
        END;
        -- delete extra user
        DELETE FROM users u WHERE u.id = usr.id[2];
        -- update merged user
        UPDATE users u SET
            email = usr.email,
            updated_date = NOW(),
            country = COALESCE(usr.country[1], usr.country[2]),
            company = COALESCE(usr.company[1], usr.company[2]),
            job_title = COALESCE(usr.job_title[1], usr.job_title[2]),
            locale = COALESCE(usr.locale[1], usr.locale[2])
            WHERE u.id = usr.id[1];

        name := usr.name[1];
        email := usr.email;

        RETURN NEXT;
    END LOOP;

    -- update active_countries
    REFRESH MATERIALIZED VIEW active_countries;
END;
$$ LANGUAGE plpgsql;

-- Clean up Guest Users (and their created battles) older than X Days --
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE last_active < (NOW() - daysOld * interval '1 day') AND type = 'GUEST';
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Clean up Battles older than X Days --
CREATE OR REPLACE PROCEDURE clean_battles(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;