CREATE TABLE IF NOT EXISTS "battles_leaders" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    PRIMARY KEY ("battle_id","user_id")
);