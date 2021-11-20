CREATE TABLE IF NOT EXISTS "api_keys" (
    "id" text NOT NULL,
    "user_id" uuid NOT NULL,
    "name" varchar(256) NOT NULL,
    "active" bool DEFAULT true,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);