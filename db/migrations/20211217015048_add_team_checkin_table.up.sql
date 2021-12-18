CREATE TABLE team_checkin (
    "id" uuid NOT NULL DEFAULT uuid_generate_v4(),
    "team_id" uuid NOT NULL REFERENCES "team" ("id") ON DELETE CASCADE,
    "user_id" uuid NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE,
    "yesterday" text,
    "today" text,
    "blockers" text,
    "discuss" text,
    "goals_met" bool DEFAULT true,
    "created_date" timestamptz NOT NULL DEFAULT now(),
    "updated_date" timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

CREATE FUNCTION prune_team_checkins()
RETURNS trigger AS $$
DECLARE
  row_count int;
BEGIN
  DELETE FROM team_checkin WHERE created_date < (NOW() - '60 days'::interval); -- clean up old checkins
  IF found THEN
    GET DIAGNOSTICS row_count = ROW_COUNT;
    RAISE NOTICE 'DELETED % row(s) FROM team_checkin', row_count;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER prune_team_checkins AFTER INSERT ON team_checkin EXECUTE PROCEDURE prune_team_checkins();