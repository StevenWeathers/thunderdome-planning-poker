ALTER TABLE ONLY "users" ALTER COLUMN "type" SET DEFAULT 'GUEST';
UPDATE "users" SET "type" = 'GUEST' WHERE "type" = 'PRIVATE';
UPDATE "users" SET "type" = 'REGISTERED' WHERE "type" = 'CORPORAL';
UPDATE "users" SET "type" = 'ADMIN' WHERE "type" = 'GENERAL';