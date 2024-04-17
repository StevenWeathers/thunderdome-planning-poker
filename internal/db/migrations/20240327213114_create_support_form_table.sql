-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.support (
    support_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    user_name character varying(128),
    user_email character varying(320),
    user_question TEXT,
    resolved boolean DEFAULT false,
    resolved_by uuid REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE thunderdome.support;
-- +goose StatementEnd
