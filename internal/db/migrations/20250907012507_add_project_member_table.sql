-- +goose Up
-- +goose StatementBegin
-- Project user table linking users to projects with role
CREATE TABLE IF NOT EXISTS thunderdome.project_user (
	project_id UUID NOT NULL REFERENCES thunderdome.project(id) ON DELETE CASCADE,
	user_id UUID NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
	role VARCHAR(16) NOT NULL DEFAULT 'MEMBER', -- ADMIN | MEMBER
	created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (project_id, user_id),
	CONSTRAINT project_user_role_check CHECK (role IN ('ADMIN','MEMBER'))
);

-- Index to optimize lookups
CREATE INDEX IF NOT EXISTS project_user_user_id_idx ON thunderdome.project_user(user_id);
CREATE INDEX IF NOT EXISTS project_user_project_id_idx ON thunderdome.project_user(project_id);

-- Project user invite table (mirrors team/department/org invite patterns)
CREATE TABLE IF NOT EXISTS thunderdome.project_user_invite (
		invite_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
		project_id uuid NOT NULL REFERENCES thunderdome.project(id) ON DELETE CASCADE,
		email character varying(320) NOT NULL,
		role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
		created_date timestamp with time zone DEFAULT now(),
		expire_date timestamp with time zone DEFAULT (now() + '24:00:00'::interval)
);

-- Prune expired project user invites function & trigger
CREATE OR REPLACE FUNCTION thunderdome.prune_project_user_invites() RETURNS trigger
		LANGUAGE plpgsql
AS $$
DECLARE
	row_count int;
BEGIN
	DELETE FROM thunderdome.project_user_invite WHERE expire_date < NOW();
	IF found THEN
		GET DIAGNOSTICS row_count = ROW_COUNT;
		RAISE NOTICE 'DELETED % row(s) FROM project_user_invite', row_count;
	END IF;
	RETURN NULL;
END;
$$;
CREATE TRIGGER prune_project_user_invites AFTER INSERT ON thunderdome.project_user_invite FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_project_user_invites();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS thunderdome.project_user_user_id_idx;
DROP INDEX IF EXISTS thunderdome.project_user_project_id_idx;
DROP TABLE IF EXISTS thunderdome.project_user;
DROP TRIGGER IF EXISTS prune_project_user_invites ON thunderdome.project_user_invite;
DROP FUNCTION IF EXISTS thunderdome.prune_project_user_invites();
DROP TABLE IF EXISTS thunderdome.project_user_invite;
-- +goose StatementEnd
