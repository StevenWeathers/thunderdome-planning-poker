-- Update User Password --
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;