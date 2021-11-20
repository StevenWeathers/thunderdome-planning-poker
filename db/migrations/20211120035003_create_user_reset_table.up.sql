CREATE TABLE IF NOT EXISTS "user_reset" (
    "reset_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '01:00:00'::interval),
    PRIMARY KEY ("reset_id")
);