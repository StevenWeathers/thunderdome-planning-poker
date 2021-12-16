-- gen_random_uuid() is faster than uuid-ossp uuid_generate_v4()
ALTER TABLE ONLY alert ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY battles ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY organization ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY organization_department ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY plans ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY team ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY user_reset ALTER COLUMN reset_id SET DEFAULT gen_random_uuid();
ALTER TABLE ONLY user_verify ALTER COLUMN verify_id SET DEFAULT gen_random_uuid();