-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.users ADD COLUMN provider TEXT NOT NULL DEFAULT 'internal';
ALTER TABLE thunderdome.users ADD COLUMN picture_url TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS provider_email_unique_idx ON thunderdome.users USING btree (provider,lower((email)::text));
DROP INDEX thunderdome.email_unique_idx;
CREATE TABLE IF NOT EXISTS thunderdome.auth_nonce (
    nonce_id character varying(64) NOT NULL PRIMARY KEY,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '10 minutes'::interval)
);
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
DROP INDEX thunderdome.provider_email_unique_idx;
ALTER TABLE thunderdome.users DROP COLUMN provider;
ALTER TABLE thunderdome.users DROP COLUMN picture_url;
DROP TABLE thunderdome.auth_nonce;
DROP TRIGGER prune_auth_nonces ON thunderdome.auth_nonce;
DROP FUNCTION thunderdome.prune_auth_nonces();
-- +goose StatementEnd
