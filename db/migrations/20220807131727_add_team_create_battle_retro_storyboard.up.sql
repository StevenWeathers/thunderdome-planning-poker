CREATE FUNCTION team_create_battle(teamid uuid, leaderid uuid, battlename character varying, pointsallowed jsonb, autovoting boolean, pointaveragerounding character varying, joincode text, leadercode text, OUT battleid uuid) RETURNS uuid
	language plpgsql
as $$
BEGIN
    INSERT INTO battles (owner_id, name, point_values_allowed, auto_finish_voting, point_average_rounding, join_code, leader_code)
        VALUES (leaderId, battleName, pointsAllowed, autoVoting, pointAverageRounding, joinCode, leaderCode)
        RETURNING id INTO battleId;
    INSERT INTO battles_leaders (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO battles_users (battle_id, user_id) VALUES (battleId, leaderId);
    INSERT INTO team_battle (team_id, battle_id) VALUES (teamid, battleId);
END;
$$;

CREATE FUNCTION team_create_retro(teamid uuid, userid uuid, retroname character varying, fmt character varying, joincode text, facilitatorcode text, maxvotes smallint, brainstormvisibility character varying) RETURNS uuid
	language plpgsql
as $$
DECLARE retroId UUID;
BEGIN
    INSERT INTO retro (owner_id, name, format, join_code, facilitator_code, max_votes, brainstorm_visibility)
    VALUES (userId, retroName, fmt, joinCode, facilitatorCode, maxVotes, brainstormVisibility) RETURNING id INTO retroId;
    INSERT INTO retro_facilitator (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO retro_user (retro_id, user_id) VALUES (retroId, userId);
    INSERT INTO team_retro (team_id, retro_id) VALUES (teamid, retroId);

    RETURN retroId;
END;
$$;

CREATE FUNCTION team_create_storyboard(teamid uuid, ownerid uuid, storyboardname character varying, joincode text, facilitatorcode text) RETURNS uuid
	language plpgsql
as $$
DECLARE storyId UUID;
BEGIN
    INSERT INTO storyboard (owner_id, name, join_code, facilitator_code)
        VALUES (ownerId, storyboardName, joinCode, facilitatorCode) RETURNING id INTO storyId;
    INSERT INTO storyboard_facilitator (storyboard_id, user_id) VALUES (storyId, ownerId);
    INSERT INTO storyboard_user (storyboard_id, user_id) VALUES(storyId, ownerId);
    INSERT INTO team_storyboard (team_id, storyboard_id) VALUES (teamid, storyId);

    RETURN storyId;
END;
$$;