ALTER TABLE users ADD COLUMN mfa_enabled bool NOT NULL DEFAULT false;
ALTER TABLE user_session ADD COLUMN disabled BOOL NOT NULL DEFAULT false;

CREATE TABLE user_mfa (
    user_id uuid PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    secret text NOT NULL,
    created_date timestamptz DEFAULT now()
);

CREATE PROCEDURE user_mfa_enable(userId UUID, mfaSecret text)
LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO user_mfa (user_id, secret) VALUES (userId, mfaSecret);
    UPDATE users SET mfa_enabled = true, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

CREATE PROCEDURE user_mfa_remove(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM user_mfa WHERE user_id = userId;
    UPDATE users SET mfa_enabled = false, updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;