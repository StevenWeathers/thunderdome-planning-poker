-- +goose Up
-- +goose StatementBegin
-- if schema_migrations exists and is not version raise error 20230820032023
DO $$
DECLARE
  v2_migrations_exist bool;
  v2_migrations_version int8;
BEGIN
    SELECT EXISTS (
        SELECT FROM
            information_schema.tables
        WHERE
            table_schema LIKE 'public' AND
            table_type LIKE 'BASE TABLE' AND
            table_name = 'schema_migrations'
    ) INTO v2_migrations_exist;
  IF v2_migrations_exist THEN
    SELECT version FROM public.schema_migrations INTO v2_migrations_version;
    IF v2_migrations_version <> 20230820032023 THEN
        RAISE EXCEPTION 'Please run Thunderdome v2.41.0 to finish v2 migrations before running v3';
    END IF;
  END IF;
END $$;
DROP TABLE IF EXISTS public.schema_migrations;
CREATE SCHEMA IF NOT EXISTS thunderdome;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA thunderdome;
-- +goose StatementEnd
