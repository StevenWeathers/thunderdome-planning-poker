-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE thunderdome.organization_user_remove(IN orgid uuid, IN userid uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM thunderdome.organization_department WHERE organization_id = orgId
    LOOP
        CALL thunderdome.department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT t.id
        FROM thunderdome.team t
        WHERE t.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.organization_user WHERE organization_id = orgId AND user_id = userId;

    COMMIT;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE thunderdome.organization_user_remove(IN orgid uuid, IN userid uuid)
    LANGUAGE plpgsql
    AS $$
DECLARE temprow record;
BEGIN
    FOR temprow IN
        SELECT id FROM thunderdome.organization_department WHERE organization_id = orgId
    LOOP
        CALL department_user_remove(temprow.id, userId);
    END LOOP;
    DELETE FROM thunderdome.team_user tu WHERE tu.team_id IN (
        SELECT t.id
        FROM thunderdome.team t
        WHERE t.organization_id = orgId
    ) AND tu.user_id = userId;
    DELETE FROM thunderdome.organization_user WHERE organization_id = orgId AND user_id = userId;

    COMMIT;
END;
$$;
-- +goose StatementEnd
