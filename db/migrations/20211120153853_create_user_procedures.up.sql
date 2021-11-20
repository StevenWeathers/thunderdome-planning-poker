-- Deletes a user and all his battle(s), api keys --
CREATE OR REPLACE PROCEDURE delete_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    DELETE FROM users WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;

    COMMIT;
END;
$$;

-- Demote User to registered type by ID --
CREATE OR REPLACE PROCEDURE demote_user(userId UUID)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET type = 'REGISTERED', updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

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

-- Reset User Password --
CREATE OR REPLACE PROCEDURE reset_user_password(resetId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM user_reset wr
        LEFT JOIN users w ON w.id = wr.user_id
        WHERE wr.reset_id = resetId AND NOW() < wr.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete in case reset record expired
        DELETE FROM user_reset WHERE reset_id = resetId;
        RAISE 'Valid Reset ID not found';
    END IF;

    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_reset WHERE reset_id = resetId;

    COMMIT;
END;
$$;

-- Update User Password --
CREATE OR REPLACE PROCEDURE update_user_password(userId UUID, userPassword TEXT)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users SET password = userPassword, last_active = NOW(), updated_date = NOW() WHERE id = userId;

    COMMIT;
END;
$$;

-- Updates a users profile --
CREATE OR REPLACE PROCEDURE user_profile_update(
    userId UUID,
    userName VARCHAR(64),
    userAvatar VARCHAR(128),
    notificationsEnabled BOOLEAN,
    userCountry VARCHAR(2),
    userLocale VARCHAR(2),
    userCompany VARCHAR(256),
    userJobTitle VARCHAR(128)
)
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE users
    SET
        name = userName,
        avatar = userAvatar,
        notifications_enabled = notificationsEnabled,
        country = userCountry,
        locale = userLocale,
        company = userCompany,
        job_title = userJobTitle,
        last_active = NOW(),
        updated_date = NOW()
    WHERE id = userId;
    REFRESH MATERIALIZED VIEW active_countries;
END;
$$;

-- Verify a user account email
CREATE OR REPLACE PROCEDURE verify_user_account(verifyId UUID)
LANGUAGE plpgsql AS $$
DECLARE matchedUserId UUID;
BEGIN
	matchedUserId := (
        SELECT w.id
        FROM user_verify wv
        LEFT JOIN users w ON w.id = wv.user_id
        WHERE wv.verify_id = verifyId AND NOW() < wv.expire_date
    );

    IF matchedUserId IS NULL THEN
        -- attempt delete incase verify record expired
        DELETE FROM user_verify WHERE verify_id = verifyId;
        RAISE 'Valid Verify ID not found';
    END IF;

    UPDATE users SET verified = 'TRUE', last_active = NOW(), updated_date = NOW() WHERE id = matchedUserId;
    DELETE FROM user_verify WHERE verify_id = verifyId;

    COMMIT;
END;
$$;