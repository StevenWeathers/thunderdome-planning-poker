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

CREATE TABLE IF NOT EXISTS "api_keys" (
    "id" text NOT NULL,
    "user_id" uuid NOT NULL,
    "name" varchar(256) NOT NULL,
    "active" bool DEFAULT true,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "battles" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "owner_id" uuid,
    "name" varchar(256),
    "voting_locked" bool DEFAULT true,
    "active_plan_id" uuid,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "point_values_allowed" jsonb DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::jsonb,
    "auto_finish_voting" bool DEFAULT true,
    "point_average_rounding" varchar(5) DEFAULT 'ceil',
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "battles_leaders" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    PRIMARY KEY ("battle_id","user_id")
);

CREATE TABLE IF NOT EXISTS "battles_users" (
    "battle_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "active" bool DEFAULT false,
    "abandoned" bool DEFAULT false,
    "spectator" bool DEFAULT false,
    PRIMARY KEY ("battle_id","user_id")
);

CREATE TABLE IF NOT EXISTS "department_team" (
    "department_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","team_id")
);

CREATE TABLE IF NOT EXISTS "department_user" (
    "department_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("department_id","user_id")
);

CREATE TABLE IF NOT EXISTS "organization" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "organization_department" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "organization_id" uuid,
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "organization_team" (
    "organization_id" uuid NOT NULL,
    "team_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","team_id")
);

CREATE TABLE IF NOT EXISTS "organization_user" (
    "organization_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("organization_id","user_id")
);

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

CREATE TABLE IF NOT EXISTS "team" (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "name" varchar(256),
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "team_battle" (
    "team_id" uuid NOT NULL,
    "battle_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    PRIMARY KEY ("team_id","battle_id")
);

CREATE TABLE IF NOT EXISTS "team_user" (
    "team_id" uuid NOT NULL,
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "updated_date" timestamp DEFAULT now(),
    "role" varchar(16) NOT NULL DEFAULT 'MEMBER',
    PRIMARY KEY ("team_id","user_id")
);

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

CREATE TABLE IF NOT EXISTS "user_reset" (
    "reset_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '01:00:00'::interval),
    PRIMARY KEY ("reset_id")
);

CREATE TABLE IF NOT EXISTS "user_verify" (
    "verify_id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "user_id" uuid NOT NULL,
    "created_date" timestamp DEFAULT now(),
    "expire_date" timestamp DEFAULT (now() + '24:00:00'::interval),
    PRIMARY KEY ("verify_id")
);