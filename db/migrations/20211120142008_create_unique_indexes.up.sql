CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx
    ON users (LOWER(email));
CREATE UNIQUE INDEX IF NOT EXISTS api_keys_warrior_id_name_key
	ON api_keys (user_id, name);
CREATE UNIQUE INDEX IF NOT EXISTS organization_team_team_id_key
	ON organization_team (team_id);
CREATE UNIQUE INDEX IF NOT EXISTS organization_department_organization_id_name_key
    ON organization_department (organization_id, name);
CREATE UNIQUE INDEX IF NOT EXISTS department_team_team_id_key
	ON department_team (team_id);
