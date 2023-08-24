CREATE OR REPLACE FUNCTION thunderdome.user_merge_nonunique_accounts()
 RETURNS TABLE(name character varying, email character varying)
 LANGUAGE plpgsql
AS $function$
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
        FROM thunderdome.users su
        WHERE su.email IS NOT NULL
        GROUP BY lower(su.email) HAVING count(su.*) > 1
        ORDER BY active_date DESC
    LOOP
        -- update poker
        UPDATE thunderdome.poker SET owner_id = usr.id[1] WHERE owner_id = usr.id[2];
        -- update poker_user
        BEGIN
            UPDATE thunderdome.poker_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in poker game';
        END;
        -- update poker_facilitator
        BEGIN
            UPDATE thunderdome.poker_facilitator SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already a poker game facilitator';
        END;
        -- update organization_user
        BEGIN
            UPDATE thunderdome.organization_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in organization';
        END;
        -- update department_user
        BEGIN
            UPDATE thunderdome.department_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in department';
        END;
        -- update team_user
        BEGIN
            UPDATE thunderdome.team_user SET user_id = usr.id[1] WHERE user_id = usr.id[2];
            EXCEPTION WHEN unique_violation THEN
            RAISE NOTICE 'User already in team';
        END;
        -- delete extra user
        DELETE FROM thunderdome.users u WHERE u.id = usr.id[2];
        -- update merged user
        UPDATE thunderdome.users u SET
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
END$function$;
