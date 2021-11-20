CREATE TABLE IF NOT EXISTS "organization_department" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "organization_id" uuid,
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);