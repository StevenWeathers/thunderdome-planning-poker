-- Add Battle Leaders by Emails --
CREATE OR REPLACE FUNCTION add_battle_leaders_by_email(
    IN battleId UUID,
    IN leaderEmails TEXT,
    OUT leaders JSONB
) AS $$
DECLARE
    emails TEXT[];
    leaderEmail TEXT;
BEGIN
    select into emails regexp_split_to_array(leaderEmails,',');
    FOREACH leaderEmail IN ARRAY emails
    LOOP
        INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, (
            SELECT id FROM users WHERE email = leaderEmail
        ));
    END LOOP;

    SELECT CASE WHEN COUNT(bl) = 0 THEN '[]'::json ELSE array_to_json(array_agg(bl.user_id)) END
    FROM battles_leaders bl WHERE bl.battle_id = battleId INTO leaders;
END;
$$ LANGUAGE plpgsql;