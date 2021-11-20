CREATE TABLE IF NOT EXISTS "organization_user" (
    "organization_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","user_id")
);