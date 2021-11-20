-- Promote User to ADMIN type by ID --
CREATE OR REPLACE PROCEDURE promote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'ADMIN', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Promote User to ADMIN type by Email --
CREATE OR REPLACE PROCEDURE promote_user_by_email(userEmail VARCHAR(320))
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'ADMIN', updated_date = NOW() WHERE email = userEmail;

    COMMIT;
END;
$$;