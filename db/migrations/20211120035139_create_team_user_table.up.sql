CREATE TABLE IF NOT EXISTS "team_user" (
    "team_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    PRIMARY KEY ("team_id","user_id")
);