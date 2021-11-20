CREATE TABLE IF NOT EXISTS "department_team" (
    "department_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","team_id")
);