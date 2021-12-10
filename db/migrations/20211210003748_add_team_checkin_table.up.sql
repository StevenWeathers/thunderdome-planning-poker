CREATE TABLE team_checkin (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "team_id" uuid NOT NULL REFERENCES "team" ("id") ON DELETE CASCADE,
    "user_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "yesterday" text,
    "today" text,
    "blockers" text,
    "discuss" text,
    "goals_met" bool DEFAULT true,
    "created_date" timestamptz NOT NULL DEFAULT now(),
    "updated_date" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);