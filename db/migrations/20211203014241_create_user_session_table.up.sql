CREATE TABLE IF NOT EXISTS "user_session" (
    "session_id" varchar(64) NOT NULL,
    "user_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '7 days'::interval),
    PRIMARY KEY ("session_id")
);