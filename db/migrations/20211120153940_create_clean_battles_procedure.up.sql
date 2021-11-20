-- Clean up Battles older than X Days --
CREATE OR REPLACE PROCEDURE clean_battles(daysOld INTEGER)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM battles WHERE updated_date < (NOW() - daysOld * interval '1 day');

    COMMIT;
END;
$$;