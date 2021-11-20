ALTER TABLE ONLY "users" ALTER COLUMN "type" SET DEFAULT 'PRIVATE';
UPDATE "users" SET "type" = 'PRIVATE' WHERE "type" = 'GUEST';
UPDATE "users" SET "type" = 'CORPORAL' WHERE "type" = 'REGISTERED';
UPDATE "users" SET "type" = 'GENERAL' WHERE "type" = 'ADMIN';