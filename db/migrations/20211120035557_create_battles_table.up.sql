CREATE TABLE IF NOT EXISTS "battles" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "owner_id" uuid,
    "name" varchar(256),
    "voting_locked" bool DEFAULT true,
    "active_plan_id" uuid,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "point_values_allowed" jsonb DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::jsonb,
    "auto_finish_voting" bool DEFAULT true,
    "point_average_rounding" varchar(5) DEFAULT 'ceil',
    PRIMARY KEY ("id")
);