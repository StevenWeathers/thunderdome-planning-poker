-- Delete Team --
DROP PROCEDURE team_delete(UUID);
CREATE PROCEDURE team_delete(teamId UUID)
AS $$
BEGIN
    DELETE FROM team WHERE id = teamId;
END;
$$ LANGUAGE plpgsql;

DROP PROCEDURE department_delete(UUID);
DROP PROCEDURE organization_delete(UUID);