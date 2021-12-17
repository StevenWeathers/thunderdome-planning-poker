DROP PROCEDURE user_account_update(
    userId UUID,
    userName VARCHAR(64),
    userEmail VARCHAR(320),
    userAvatar VARCHAR(128),
    notificationsEnabled BOOLEAN,
    userCountry VARCHAR(2),
    userLocale VARCHAR(2),
    userCompany VARCHAR(256),
    userJobTitle VARCHAR(128)
);