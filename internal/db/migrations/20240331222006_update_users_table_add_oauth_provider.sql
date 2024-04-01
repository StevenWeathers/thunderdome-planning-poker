-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.users ADD COLUMN provider TEXT NOT NULL DEFAULT 'internal';
ALTER TABLE thunderdome.users ADD COLUMN picture_url TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS provider_email_unique_idx ON thunderdome.users USING btree (provider,lower((email)::text));
DROP INDEX thunderdome.email_unique_idx;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx ON thunderdome.users USING btree (lower((email)::text));
DROP INDEX thunderdome.provider_email_unique_idx;
ALTER TABLE thunderdome.users DROP COLUMN provider;
ALTER TABLE thunderdome.users DROP COLUMN picture_url;
-- +goose StatementEnd
