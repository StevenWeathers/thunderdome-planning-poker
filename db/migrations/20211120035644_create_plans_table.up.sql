CREATE TABLE IF NOT EXISTS "plans" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "points" varchar(3) DEFAULT '',
    "active" bool DEFAULT false,
    "battle_id" uuid NOT NULL,
    "votes" jsonb DEFAULT '[]'::jsonb,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "skipped" bool DEFAULT false,
    "votestart_time" timestamp DEFAULT now(),
    "voteend_time" timestamp DEFAULT now(),
    "acceptance_criteria" text,
    "link" text,
    "description" text,
    "reference_id" varchar(128),
    "type" varchar(64) DEFAULT 'story',
    PRIMARY KEY ("id")
);