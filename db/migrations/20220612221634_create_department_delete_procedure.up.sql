-- Delete Team --
DROP PROCEDURE team_delete(UUID);
CREATE PROCEDURE team_delete(teamId UUID)
AS $$
BEGIN
    DELETE FROM battles WHERE id IN (
        SELECT battle_id FROM team_battle WHERE team_id = teamId
    );

    DELETE FROM retro WHERE id IN (
        SELECT retro_id FROM team_retro WHERE team_id = teamId
    );

    DELETE FROM storyboard WHERE id IN (
        SELECT storyboard_id FROM team_storyboard WHERE team_id = teamId
    );

    DELETE FROM team WHERE id = teamId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

-- Delete Department --
CREATE PROCEDURE department_delete(deptId UUID)
AS $$
    DECLARE t record;
BEGIN
    FOR t IN SELECT team_id FROM department_team WHERE department_id = deptId
    LOOP
	    CALL team_delete(t.team_id);
    END LOOP;

    DELETE FROM organization_department WHERE id = deptId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;

-- Delete Organization --
CREATE PROCEDURE organization_delete(orgId UUID)
AS $$
    DECLARE d record;
BEGIN
    FOR d IN SELECT id FROM organization_department WHERE organization_id = orgId
    LOOP
	    CALL department_delete(d.id);
    END LOOP;

    DELETE FROM organization WHERE id = orgId;

    COMMIT;
END;
$$ LANGUAGE plpgsql;