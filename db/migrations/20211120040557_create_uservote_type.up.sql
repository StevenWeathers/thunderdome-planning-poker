DROP TYPE IF EXISTS UsersVote;
CREATE TYPE UsersVote AS
(
    "warriorId"     uuid,
    "vote"   VARCHAR(3)
);