CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx
    ON users (LOWER(email));
ALTER TABLE api_keys DROP CONSTRAINT IF EXISTS api_keys_warrior_id_name_key;
ALTER TABLE api_keys ADD CONSTRAINT api_keys_user_id_name_key
	UNIQUE (user_id, name);
ALTER TABLE organization_team DROP CONSTRAINT IF EXISTS organization_team_team_id_key;
ALTER TABLE organization_team ADD CONSTRAINT organization_team_team_id_key
    UNIQUE (team_id);
ALTER TABLE organization_department DROP CONSTRAINT IF EXISTS organization_department_organization_id_name_key;
ALTER TABLE organization_department ADD CONSTRAINT organization_department_organization_id_name_key
    UNIQUE (organization_id, name);
ALTER TABLE department_team DROP CONSTRAINT IF EXISTS department_team_team_id_key;
ALTER TABLE department_team ADD CONSTRAINT department_team_team_id_key
    UNIQUE (team_id);
