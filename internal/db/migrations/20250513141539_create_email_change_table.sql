-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.user_email_change (
    change_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '01:00:00'::interval)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS thunderdome.user_email_change;
-- +goose StatementEnd
