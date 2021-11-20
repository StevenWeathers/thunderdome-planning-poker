CREATE TABLE IF NOT EXISTS "department_user" (
    "department_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","user_id")
);