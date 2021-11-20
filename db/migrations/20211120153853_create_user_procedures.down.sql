DROP PROCEDURE delete_user(userId UUID);
DROP PROCEDURE demote_user(userId UUID);
DROP PROCEDURE promote_user(userId UUID);
DROP PROCEDURE promote_user_by_email(userEmail VARCHAR(320));
DROP PROCEDURE reset_user_password(resetId UUID, userPassword TEXT);
DROP PROCEDURE update_user_password(userId UUID, userPassword TEXT);
DROP PROCEDURE user_profile_update(
    userId UUID,
    userName VARCHAR(64),
    userAvatar VARCHAR(128),
    notificationsEnabled BOOLEAN,
    userCountry VARCHAR(2),
    userLocale VARCHAR(2),
    userCompany VARCHAR(256),
    userJobTitle VARCHAR(128)
);
DROP PROCEDURE verify_user_account(verifyId UUID);