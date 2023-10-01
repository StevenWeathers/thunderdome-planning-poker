-- +goose Up
-- +goose StatementBegin
CREATE TABLE thunderdome.jira_instance (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    host text NOT NULL,
    client_mail text NOT NULL,
    access_token text NOT NULL,
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now()
);
ALTER TABLE thunderdome.users ADD COLUMN subscribed boolean NOT NULL DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE thunderdome.jira_instance;
ALTER TABLE thunderdome.users DROP COLUMN subscribed;
-- +goose StatementEnd
