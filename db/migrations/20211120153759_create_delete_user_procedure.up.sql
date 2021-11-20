-- Deletes a user and all his battle(s), api keys --
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;