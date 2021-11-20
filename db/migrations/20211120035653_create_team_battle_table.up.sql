CREATE TABLE IF NOT EXISTS "team_battle" (
    "team_id" uuid NOT NULL,
    "battle_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("team_id","battle_id")
);