CREATE TABLE IF NOT EXISTS "battles_users" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "active" bool DEFAULT false,
    "abandoned" bool DEFAULT false,
    "spectator" bool DEFAULT false,
    PRIMARY KEY ("battle_id","user_id")
);