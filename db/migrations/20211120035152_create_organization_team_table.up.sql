CREATE TABLE IF NOT EXISTS "organization_team" (
    "organization_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","team_id")
);