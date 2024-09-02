-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS thunderdome.retro_item_comment (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    item_id uuid REFERENCES thunderdome.retro_item(id) ON DELETE CASCADE,
    user_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    comment text,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS thunderdome.retro_item_comment;
-- +goose StatementEnd
