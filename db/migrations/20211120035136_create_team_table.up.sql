CREATE TABLE IF NOT EXISTS "team" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);