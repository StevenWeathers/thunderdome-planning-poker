-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.auth_nonce (
    nonce_id character varying(64) NOT NULL PRIMARY KEY,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '10 minutes'::interval)
);
CREATE TABLE thunderdome.auth_credential (
    user_id uuid UNIQUE NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    email character varying(320),
    password text,
    verified boolean DEFAULT false,
    mfa_enabled boolean DEFAULT false NOT NULL,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX cred_email_unique_idx ON thunderdome.auth_credential USING btree (lower((email)::text));
CREATE TABLE thunderdome.auth_identity (
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    provider character varying(64) NOT NULL,
    sub TEXT NOT NULL,
    email character varying(320) NOT NULL,
    verified boolean NOT NULL DEFAULT false,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now(),
    UNIQUE(provider, sub)
);
ALTER TABLE thunderdome.users ADD COLUMN picture_url TEXT;
CREATE OR REPLACE FUNCTION thunderdome.prune_auth_nonces() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM thunderdome.auth_nonce WHERE expire_date < NOW();
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM auth_nonce', row_count;
  END IF;
  RETURN NULL;
END;
$$;
CREATE TRIGGER prune_auth_nonces AFTER INSERT ON thunderdome.auth_nonce FOR EACH STATEMENT EXECUTE FUNCTION thunderdome.prune_auth_nonces();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx ON thunderdome.users USING btree (lower((email)::text));
DROP TRIGGER prune_auth_nonces ON thunderdome.auth_nonce;
DROP FUNCTION thunderdome.prune_auth_nonces();
ALTER TABLE thunderdome.users DROP COLUMN picture_url;
DROP TABLE thunderdome.auth_nonce;
DROP TABLE thunderdome.auth_credential;
DROP TABLE thunderdome.auth_identity;
-- +goose StatementEnd
