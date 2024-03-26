-- +goose Up
-- +goose StatementBegin
-- setting temp default subscription ID for existing single user subs will need to be updated immediately
ALTER TABLE thunderdome.subscription ADD COLUMN subscription_id TEXT NOT NULL DEFAULT 'temp_USERSUB';
ALTER TABLE thunderdome.subscription ADD COLUMN type TEXT NOT NULL DEFAULT 'user';
ALTER TABLE thunderdome.subscription DROP CONSTRAINT subscription_user_id_key;
CREATE UNIQUE INDEX IF NOT EXISTS subscription_user_id_customer_id_idx ON thunderdome.subscription USING btree (user_id, customer_id, subscription_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE thunderdome.subscription DROP COLUMN subscription_id;
ALTER TABLE thunderdome.subscription DROP COLUMN type;
DROP INDEX thunderdome.subscription_user_id_customer_id_idx;
CREATE UNIQUE INDEX IF NOT EXISTS subscription_user_id_key ON thunderdome.subscription USING btree (user_id);
-- +goose StatementEnd
