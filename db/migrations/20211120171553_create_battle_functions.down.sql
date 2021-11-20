DROP FUNCTION create_battle(
    IN leaderId UUID,
    IN battleName VARCHAR(256),
    IN pointsAllowed JSONB,
    IN autoVoting BOOL,
    IN pointAverageRounding VARCHAR(5),
    OUT battleId UUID
);

DROP FUNCTION add_battle_leaders_by_email(
    IN battleId UUID,
    IN leaderEmails TEXT,
    OUT leaders JSONB
);