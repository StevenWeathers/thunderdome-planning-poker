DROP FUNCTION user_session_get(
    IN sessionId VARCHAR(64)
);
DROP TRIGGER prune_user_sessions ON user_session;
DROP FUNCTION prune_user_sessions();