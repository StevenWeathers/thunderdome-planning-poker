-- +goose Up
ALTER TABLE thunderdome.jira_instance ADD COLUMN jira_data_center boolean NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE thunderdome.jira_instance DROP COLUMN jira_data_center;

