-- Demote User to registered type by ID --
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'REGISTERED', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;