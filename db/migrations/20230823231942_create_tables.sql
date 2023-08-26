-- +goose Up
-- +goose StatementBegin
DROP TYPE IF EXISTS thunderdome.usersvote;
CREATE TYPE thunderdome.usersvote AS (
    "warriorId" uuid,
    vote character varying(3)
);

CREATE TABLE IF NOT EXISTS thunderdome.users (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(64),
    created_date timestamp with time zone DEFAULT now(),
    last_active timestamp with time zone DEFAULT now(),
    email character varying(320),
    password text,
    type character varying(128) DEFAULT 'GUEST'::character varying,
    verified boolean DEFAULT false,
    avatar character varying(128) DEFAULT 'robohash'::character varying,
    notifications_enabled boolean DEFAULT true,
    country character varying(2),
    company character varying(256),
    job_title character varying(128),
    updated_date timestamp with time zone DEFAULT now(),
    locale character varying(2),
    disabled boolean DEFAULT false,
    mfa_enabled boolean DEFAULT false NOT NULL,
    theme character varying(5) DEFAULT 'auto'::character varying NOT NULL
);

CREATE MATERIALIZED VIEW IF NOT EXISTS thunderdome.active_countries AS
    SELECT DISTINCT users.country FROM thunderdome.users WITH NO DATA;

CREATE TABLE IF NOT EXISTS thunderdome.alert (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(256) NOT NULL,
    type character varying(128) DEFAULT 'NEW'::character varying,
    content text NOT NULL,
    active boolean DEFAULT true,
    allow_dismiss boolean DEFAULT true,
    registered_only boolean DEFAULT true,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.api_key (
    id text NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    name character varying(256) NOT NULL,
    active boolean DEFAULT true,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    UNIQUE (user_id, name)
);

CREATE TABLE IF NOT EXISTS thunderdome.organization (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(256),
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.organization_department (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    organization_id uuid REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    name character varying(256),
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    UNIQUE (organization_id, name)
);

CREATE TABLE IF NOT EXISTS thunderdome.department_user (
    department_id uuid NOT NULL REFERENCES thunderdome.organization_department(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (department_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.organization_user (
    organization_id uuid NOT NULL REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (organization_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.team (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(256),
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    organization_id uuid REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    department_id uuid REFERENCES thunderdome.organization_department(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thunderdome.team_user (
    team_id uuid NOT NULL REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    role character varying(16) DEFAULT 'MEMBER'::character varying NOT NULL,
    PRIMARY KEY (team_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.poker (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    owner_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    name character varying(256),
    voting_locked boolean DEFAULT true,
    active_story_id uuid,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    point_values_allowed jsonb DEFAULT '["1/2", "1", "2", "3", "5", "8", "13", "?"]'::jsonb,
    auto_finish_voting boolean DEFAULT true,
    point_average_rounding character varying(5) DEFAULT 'ceil'::character varying,
    join_code text,
    leader_code text,
    hide_voter_identity boolean DEFAULT false,
    last_active timestamp with time zone DEFAULT now() NOT NULL,
    team_id uuid REFERENCES thunderdome.team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thunderdome.poker_facilitator (
    poker_id uuid NOT NULL REFERENCES thunderdome.poker(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    PRIMARY KEY (poker_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.poker_story (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(256),
    points character varying(3) DEFAULT ''::character varying,
    active boolean DEFAULT false,
    poker_id uuid NOT NULL REFERENCES thunderdome.poker(id) ON DELETE CASCADE,
    votes jsonb DEFAULT '[]'::jsonb,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    skipped boolean DEFAULT false,
    votestart_time timestamp with time zone DEFAULT now(),
    voteend_time timestamp with time zone DEFAULT now(),
    link text,
    description text,
    acceptance_criteria text,
    reference_id character varying(128),
    type character varying(64) DEFAULT 'story'::character varying,
    priority integer DEFAULT 99,
    "position" double precision DEFAULT 0 NOT NULL,
    UNIQUE (poker_id, "position")
);

CREATE TABLE IF NOT EXISTS thunderdome.poker_user (
    poker_id uuid NOT NULL REFERENCES thunderdome.poker(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    active boolean DEFAULT false,
    abandoned boolean DEFAULT false,
    spectator boolean DEFAULT false,
    PRIMARY KEY (poker_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.retro (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    owner_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    name character varying(256),
    format character varying(32) DEFAULT 'worked_improve_question'::character varying NOT NULL,
    phase character varying(16) DEFAULT 'intro'::character varying NOT NULL,
    join_code text,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    brainstorm_visibility character varying(12) DEFAULT 'visible'::character varying,
    max_votes smallint DEFAULT 3,
    facilitator_code text,
    last_active timestamp with time zone DEFAULT now() NOT NULL,
    team_id uuid REFERENCES thunderdome.team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_action (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    retro_id uuid REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    content text NOT NULL,
    completed boolean DEFAULT false,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_action_assignee (
    action_id uuid NOT NULL REFERENCES thunderdome.retro_action(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (action_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_action_comment (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    action_id uuid REFERENCES thunderdome.retro_action(id) ON DELETE CASCADE,
    user_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    comment text,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_facilitator (
    retro_id uuid NOT NULL REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (retro_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_group (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    retro_id uuid REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    name character varying(128),
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_group_vote (
    retro_id uuid NOT NULL REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    group_id uuid NOT NULL REFERENCES thunderdome.retro_group(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    PRIMARY KEY (retro_id, user_id, group_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_item (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    retro_id uuid REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    user_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    group_id uuid REFERENCES thunderdome.retro_group(id) ON DELETE CASCADE,
    content text NOT NULL,
    type character varying(16) NOT NULL,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.retro_user (
    retro_id uuid NOT NULL REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    active boolean DEFAULT false,
    abandoned boolean DEFAULT false,
    spectator boolean DEFAULT false,
    PRIMARY KEY (retro_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    name character varying(256),
    owner_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    join_code text,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    color_legend jsonb DEFAULT '[{"color": "gray", "legend": ""}, {"color": "red", "legend": ""}, {"color": "orange", "legend": ""}, {"color": "yellow", "legend": ""}, {"color": "green", "legend": ""}, {"color": "teal", "legend": ""}, {"color": "blue", "legend": ""}, {"color": "indigo", "legend": ""}, {"color": "purple", "legend": ""}, {"color": "pink", "legend": ""}]'::jsonb,
    facilitator_code text,
    last_active timestamp with time zone DEFAULT now() NOT NULL,
    team_id uuid REFERENCES thunderdome.team(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_facilitator (
    storyboard_id uuid NOT NULL REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (storyboard_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_persona (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    storyboard_id uuid REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    name character varying(256) NOT NULL,
    role character varying(256),
    description text,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    UNIQUE (storyboard_id, name)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_goal (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    storyboard_id uuid REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    name character varying(256),
    sort_order integer,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    UNIQUE (storyboard_id, sort_order)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_column (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    storyboard_id uuid REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    goal_id uuid REFERENCES thunderdome.storyboard_goal(id) ON DELETE CASCADE,
    name character varying(256),
    sort_order integer,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    UNIQUE (goal_id, sort_order)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_column_persona (
    column_id uuid NOT NULL REFERENCES thunderdome.storyboard_column(id) ON DELETE CASCADE,
    persona_id uuid NOT NULL REFERENCES thunderdome.storyboard_persona(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (column_id, persona_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_goal_persona (
    goal_id uuid NOT NULL REFERENCES thunderdome.storyboard_goal(id) ON DELETE CASCADE,
    persona_id uuid NOT NULL REFERENCES thunderdome.storyboard_persona(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    PRIMARY KEY (goal_id, persona_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_story (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    storyboard_id uuid REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    goal_id uuid REFERENCES thunderdome.storyboard_goal(id) ON DELETE CASCADE,
    column_id uuid REFERENCES thunderdome.storyboard_column(id) ON DELETE CASCADE,
    name character varying(256),
    color character varying(32) DEFAULT 'gray'::character varying,
    content text,
    sort_order integer,
    points integer,
    closed boolean DEFAULT false,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now(),
    annotations jsonb DEFAULT '[]'::jsonb NOT NULL,
    link text,
    UNIQUE (column_id, sort_order) DEFERRABLE
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_story_comment (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    storyboard_id uuid REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    story_id uuid REFERENCES thunderdome.storyboard_story(id) ON DELETE CASCADE,
    comment text,
    user_id uuid REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    updated_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.storyboard_user (
    storyboard_id uuid NOT NULL REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    active boolean DEFAULT false,
    abandoned boolean DEFAULT false,
    PRIMARY KEY (storyboard_id, user_id)
);

CREATE TABLE IF NOT EXISTS thunderdome.team_checkin (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    team_id uuid NOT NULL REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    yesterday text,
    today text,
    blockers text,
    discuss text,
    goals_met boolean DEFAULT true,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    updated_date timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS thunderdome.team_checkin_comment (
    id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    checkin_id uuid NOT NULL REFERENCES thunderdome.team_checkin(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    comment text,
    created_date timestamp with time zone DEFAULT now() NOT NULL,
    updated_date timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE IF NOT EXISTS thunderdome.user_mfa (
    user_id uuid NOT NULL PRIMARY KEY REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    secret text NOT NULL,
    created_date timestamp with time zone DEFAULT now()
);

CREATE TABLE IF NOT EXISTS thunderdome.user_reset (
    reset_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '01:00:00'::interval)
);

CREATE TABLE IF NOT EXISTS thunderdome.user_session (
    session_id character varying(64) NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '7 days'::interval),
    disabled boolean DEFAULT false NOT NULL
);

CREATE TABLE IF NOT EXISTS thunderdome.user_verify (
    verify_id uuid DEFAULT gen_random_uuid() NOT NULL PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    created_date timestamp with time zone DEFAULT now(),
    expire_date timestamp with time zone DEFAULT (now() + '24:00:00'::interval)
);

CREATE INDEX IF NOT EXISTS battles_leaders_user_id_idx ON thunderdome.poker_facilitator USING btree (user_id);
CREATE INDEX IF NOT EXISTS battles_owner_id_idx ON thunderdome.poker USING btree (owner_id);
CREATE INDEX IF NOT EXISTS battles_users_user_id_idx ON thunderdome.poker_user USING btree (user_id);
CREATE INDEX IF NOT EXISTS department_user_user_id_idx ON thunderdome.department_user USING btree (user_id);
CREATE UNIQUE INDEX IF NOT EXISTS email_unique_idx ON thunderdome.users USING btree (lower((email)::text));
CREATE INDEX IF NOT EXISTS organization_user_user_id_idx ON thunderdome.organization_user USING btree (user_id);
CREATE INDEX IF NOT EXISTS plans_battle_id_idx ON thunderdome.poker_story USING btree (poker_id);
CREATE INDEX IF NOT EXISTS poker_team_id_idx ON thunderdome.poker USING btree (team_id);
CREATE INDEX IF NOT EXISTS retro_action_assignee_user_id_idx ON thunderdome.retro_action_assignee USING btree (user_id);
CREATE INDEX IF NOT EXISTS retro_action_comment_action_id_idx ON thunderdome.retro_action_comment USING btree (action_id);
CREATE INDEX IF NOT EXISTS retro_action_comment_user_id_idx ON thunderdome.retro_action_comment USING btree (user_id);
CREATE INDEX IF NOT EXISTS retro_action_retro_id_idx ON thunderdome.retro_action USING btree (retro_id);
CREATE INDEX IF NOT EXISTS retro_facilitator_user_id_idx ON thunderdome.retro_facilitator USING btree (user_id);
CREATE INDEX IF NOT EXISTS retro_group_retro_id_idx ON thunderdome.retro_group USING btree (retro_id);
CREATE INDEX IF NOT EXISTS retro_group_vote_group_id_idx ON thunderdome.retro_group_vote USING btree (group_id);
CREATE INDEX IF NOT EXISTS retro_group_vote_user_id_idx ON thunderdome.retro_group_vote USING btree (user_id);
CREATE INDEX IF NOT EXISTS retro_item_group_id_idx ON thunderdome.retro_item USING btree (group_id);
CREATE INDEX IF NOT EXISTS retro_item_retro_id_idx ON thunderdome.retro_item USING btree (retro_id);
CREATE INDEX IF NOT EXISTS retro_item_user_id_idx ON thunderdome.retro_item USING btree (user_id);
CREATE INDEX IF NOT EXISTS retro_owner_id_idx ON thunderdome.retro USING btree (owner_id);
CREATE INDEX IF NOT EXISTS retro_team_id_idx ON thunderdome.retro USING btree (team_id);
CREATE INDEX IF NOT EXISTS retro_user_user_id_idx ON thunderdome.retro_user USING btree (user_id);
CREATE INDEX IF NOT EXISTS storyboard_column_persona_persona_id_idx ON thunderdome.storyboard_column_persona USING btree (persona_id);
CREATE INDEX IF NOT EXISTS storyboard_column_storyboard_id_idx ON thunderdome.storyboard_column USING btree (storyboard_id);
CREATE INDEX IF NOT EXISTS storyboard_facilitator_user_id_idx ON thunderdome.storyboard_facilitator USING btree (user_id);
CREATE INDEX IF NOT EXISTS storyboard_goal_persona_persona_id_idx ON thunderdome.storyboard_goal_persona USING btree (persona_id);
CREATE INDEX IF NOT EXISTS storyboard_owner_id_idx ON thunderdome.storyboard USING btree (owner_id);
CREATE INDEX IF NOT EXISTS storyboard_story_comment_story_id_idx ON thunderdome.storyboard_story_comment USING btree (story_id);
CREATE INDEX IF NOT EXISTS storyboard_story_comment_storyboard_id_idx ON thunderdome.storyboard_story_comment USING btree (storyboard_id);
CREATE INDEX IF NOT EXISTS storyboard_story_comment_user_id_idx ON thunderdome.storyboard_story_comment USING btree (user_id);
CREATE INDEX IF NOT EXISTS storyboard_story_goal_id_idx ON thunderdome.storyboard_story USING btree (goal_id);
CREATE INDEX IF NOT EXISTS storyboard_story_storyboard_id_idx ON thunderdome.storyboard_story USING btree (storyboard_id);
CREATE INDEX IF NOT EXISTS storyboard_team_id_idx ON thunderdome.storyboard USING btree (team_id);
CREATE INDEX IF NOT EXISTS storyboard_user_user_id_idx ON thunderdome.storyboard_user USING btree (user_id);
CREATE INDEX IF NOT EXISTS team_checkin_comment_checkin_id_idx ON thunderdome.team_checkin_comment USING btree (checkin_id);
CREATE INDEX IF NOT EXISTS team_checkin_comment_user_id_idx ON thunderdome.team_checkin_comment USING btree (user_id);
CREATE INDEX IF NOT EXISTS team_checkin_team_id_idx ON thunderdome.team_checkin USING btree (team_id);
CREATE INDEX IF NOT EXISTS team_checkin_user_id_idx ON thunderdome.team_checkin USING btree (user_id);
CREATE INDEX IF NOT EXISTS team_department_id_idx ON thunderdome.team USING btree (department_id);
CREATE INDEX IF NOT EXISTS team_organization_id_idx ON thunderdome.team USING btree (organization_id);
CREATE INDEX IF NOT EXISTS team_user_user_id_idx ON thunderdome.team_user USING btree (user_id);
CREATE INDEX IF NOT EXISTS user_reset_user_id_idx ON thunderdome.user_reset USING btree (user_id);
CREATE INDEX IF NOT EXISTS user_session_user_id_idx ON thunderdome.user_session USING btree (user_id);
CREATE INDEX IF NOT EXISTS user_verify_user_id_idx ON thunderdome.user_verify USING btree (user_id);

ALTER TABLE ONLY thunderdome.api_key DROP CONSTRAINT IF EXISTS apk_warrior_id_fkey;
ALTER TABLE ONLY thunderdome.poker DROP CONSTRAINT IF EXISTS battles_owner_id_fkey;
ALTER TABLE ONLY thunderdome.poker_facilitator DROP CONSTRAINT IF EXISTS bl_battle_id_fkey;
ALTER TABLE ONLY thunderdome.poker_facilitator DROP CONSTRAINT IF EXISTS bl_warrior_id_fkey;
ALTER TABLE ONLY thunderdome.poker_user DROP CONSTRAINT IF EXISTS bw_battle_id_fkey;
ALTER TABLE ONLY thunderdome.poker_user DROP CONSTRAINT IF EXISTS bw_warrior_id_fkey;
ALTER TABLE ONLY thunderdome.department_user DROP CONSTRAINT IF EXISTS department_user_department_id_fkey;
ALTER TABLE ONLY thunderdome.department_user DROP CONSTRAINT IF EXISTS department_user_user_id_fkey;
ALTER TABLE ONLY thunderdome.organization_department DROP CONSTRAINT IF EXISTS organization_department_organization_id_fkey;
ALTER TABLE ONLY thunderdome.organization_user DROP CONSTRAINT IF EXISTS organization_user_organization_id_fkey;
ALTER TABLE ONLY thunderdome.organization_user DROP CONSTRAINT IF EXISTS organization_user_user_id_fkey;
ALTER TABLE ONLY thunderdome.poker_story DROP CONSTRAINT IF EXISTS plans_battle_id_fkey;
ALTER TABLE ONLY thunderdome.team_user DROP CONSTRAINT IF EXISTS team_user_team_id_fkey;
ALTER TABLE ONLY thunderdome.team_user DROP CONSTRAINT IF EXISTS team_user_user_id_fkey;
ALTER TABLE ONLY thunderdome.user_reset DROP CONSTRAINT IF EXISTS wr_warrior_id_fkey;
ALTER TABLE ONLY thunderdome.user_verify DROP CONSTRAINT IF EXISTS wv_warrior_id_fkey;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE thunderdome.usersvote;
DROP TABLE thunderdome.users;
DROP MATERIALIZED VIEW thunderdome.active_countries;
DROP TABLE thunderdome.alert;
DROP TABLE thunderdome.api_key;
DROP TABLE thunderdome.department_user;
DROP TABLE thunderdome.organization;
DROP TABLE thunderdome.organization_department;
DROP TABLE thunderdome.organization_user;
DROP TABLE thunderdome.poker;
DROP TABLE thunderdome.poker_facilitator;
DROP TABLE thunderdome.poker_story;
DROP TABLE thunderdome.poker_user;
DROP TABLE thunderdome.retro;
DROP TABLE thunderdome.retro_action;
DROP TABLE thunderdome.retro_action_assignee;
DROP TABLE thunderdome.retro_action_comment;
DROP TABLE thunderdome.retro_facilitator;
DROP TABLE thunderdome.retro_group;
DROP TABLE thunderdome.retro_group_vote;
DROP TABLE thunderdome.retro_item;
DROP TABLE thunderdome.retro_user;
DROP TABLE thunderdome.storyboard;
DROP TABLE thunderdome.storyboard_column;
DROP TABLE thunderdome.storyboard_column_persona;
DROP TABLE thunderdome.storyboard_facilitator;
DROP TABLE thunderdome.storyboard_goal;
DROP TABLE thunderdome.storyboard_goal_persona;
DROP TABLE thunderdome.storyboard_persona;
DROP TABLE thunderdome.storyboard_story;
DROP TABLE thunderdome.storyboard_story_comment;
DROP TABLE thunderdome.storyboard_user;
DROP TABLE thunderdome.team;
DROP TABLE thunderdome.team_checkin;
DROP TABLE thunderdome.team_checkin_comment;
DROP TABLE thunderdome.team_user;
DROP TABLE thunderdome.user_mfa;
DROP TABLE thunderdome.user_reset;
DROP TABLE thunderdome.user_session;
DROP TABLE thunderdome.user_verify;
-- +goose StatementEnd
