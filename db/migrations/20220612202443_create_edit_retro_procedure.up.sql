-- Edit a Retro
CREATE PROCEDURE edit_retro(retroId UUID, retroName VARCHAR(256), joinCode VARCHAR(128))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE retro SET name = retroName, join_code = joinCode, updated_date = NOW()
        WHERE id = retroId;

    COMMIT;
END;
$$;