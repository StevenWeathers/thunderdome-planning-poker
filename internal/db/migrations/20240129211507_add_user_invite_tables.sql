-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.organization_user_invite (
    invite_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    organization_id uuid NOT NULL REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    email character varying(320) NOT NULL,
    role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '24:00:00'::interval)
);
CREATE TABLE IF NOT EXISTS thunderdome.team_user_invite (
    invite_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    team_id uuid NOT NULL REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    email character varying(320) NOT NULL,
    role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '24:00:00'::interval)
);
CREATE OR REPLACE FUNCTION thunderdome.prune_team_user_invites() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.team_user_invite WHERE expire_date < NOW();
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM team_user_invite', row_count;
  END IF;
  RETURN NULL;
END;
$$;
CREATE OR REPLACE FUNCTION thunderdome.prune_organization_user_invites() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.organization_user_invite WHERE expire_date < NOW();
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM organization_user_invite', row_count;
  END IF;
  RETURN NULL;
END;
$$;
CREATE TRIGGER prune_team_user_invites AFTER INSERT ON thunderdome.team_user_invite FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_team_user_invites();
CREATE TRIGGER prune_organization_user_invites AFTER INSERT ON thunderdome.organization_user_invite FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_organization_user_invites();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS prune_team_user_invites ON thunderdome.team_user_invite;
DROP TRIGGER IF EXISTS prune_organization_user_invites ON thunderdome.organization_user_invite;
DROP FUNCTION thunderdome.prune_organization_user_invites();
DROP FUNCTION thunderdome.prune_team_user_invites();
DROP TABLE IF EXISTS thunderdome.organization_user_invite;
DROP TABLE IF EXISTS thunderdome.team_user_invite;
-- +goose StatementEnd
