CREATE TABLE IF NOT EXISTS "users" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(64),
    "created_date" timestamp DEFAULT now(),
    "last_active" timestamp DEFAULT now(),
    "email" varchar(320),
    "password" text,
    "type" varchar(128) DEFAULT 'PRIVATE',
    "verified" bool DEFAULT false,
    "avatar" varchar(128) DEFAULT 'identicon',
    "notifications_enabled" bool DEFAULT true,
    "country" varchar(2),
    "company" varchar(256),
    "job_title" varchar(128),
    "updated_date" timestamp DEFAULT now(),
    "locale" varchar(2),
    PRIMARY KEY ("id")
);