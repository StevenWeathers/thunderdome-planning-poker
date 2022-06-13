CREATE TABLE team_checkin_comment (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "checkin_id" uuid NOT NULL REFERENCES "team_checkin" ("id") ON DELETE CASCADE,
    "user_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "comment" text,
    "created_date" timestamptz NOT NULL DEFAULT now(),
    "updated_date" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);