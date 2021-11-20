CREATE TABLE IF NOT EXISTS "alert" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "type" varchar(128) DEFAULT 'NEW',
    "content" text,
    "active" bool DEFAULT true,
    "allow_dismiss" bool DEFAULT true,
    "registered_only" bool DEFAULT true,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);