-- Clean up Guest Users (and their created battles) older than X Days --
CREATE OR REPLACE PROCEDURE clean_guest_users(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE last_active < (NOW() - daysOld * interval '1 day') AND type = 'GUEST';
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;