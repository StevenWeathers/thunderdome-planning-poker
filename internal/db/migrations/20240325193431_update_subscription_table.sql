-- +goose Up
-- +goose StatementBegin
ALTER TABLE thunderdome.subscription DROP CONSTRAINT subscription_user_id_key;
CREATE UNIQUE INDEX IF NOT EXISTS subscription_user_id_customer_id_idx ON thunderdome.subscription USING btree (user_id, customer_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX thunderdome.subscription_user_id_customer_id_idx;
CREATE UNIQUE INDEX IF NOT EXISTS subscription_user_id_key ON thunderdome.subscription USING btree (user_id);
-- +goose StatementEnd
