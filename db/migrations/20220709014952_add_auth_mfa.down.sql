ALTER TABLE users DROP COLUMN mfa_enabled;
DROP TABLE user_mfa;
DROP PROCEDURE user_mfa_enable(UUID, text);
DROP PROCEDURE user_mfa_remove(UUID);