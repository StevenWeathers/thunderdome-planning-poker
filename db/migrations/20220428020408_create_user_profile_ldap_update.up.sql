-- Updates a users profile (excludes username) --
CREATE OR REPLACE PROCEDURE user_profile_ldap_update(
    userId UUID,
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