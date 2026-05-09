-- +goose Up
-- +goose StatementBegin
UPDATE thunderdome.user_session
SET disabled = NOT disabled;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE thunderdome.user_session
SET disabled = NOT disabled;
-- +goose StatementEnd