-- +goose Up
-- +goose StatementBegin
CREATE TABLE thunderdome.subscription (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    customer_id text UNIQUE NOT NULL,
    active boolean NOT NULL DEFAULT true,
    expires timestamp with time zone NOT NULL DEFAULT NOW(),
    created_date timestamp with time zone NOT NULL DEFAULT now(),
    updated_date timestamp with time zone NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE thunderdome.subscription;
-- +goose StatementEnd
