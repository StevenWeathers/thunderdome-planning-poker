-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS thunderdome.health_check_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    areas JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT health_check_template_check_scope CHECK (
        (is_public = true AND num_nonnulls(organization_id, team_id) = 0) OR
        (is_public = false AND num_nonnulls(organization_id, team_id) = 1)
    )
);

CREATE INDEX IF NOT EXISTS idx_health_check_template_organization_id
    ON thunderdome.health_check_template(organization_id);

CREATE INDEX IF NOT EXISTS idx_health_check_template_team_id
    ON thunderdome.health_check_template(team_id);

-- areas are copied from the template so that later template edits don't change existing health checks
CREATE TABLE IF NOT EXISTS thunderdome.health_check (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    team_id UUID NOT NULL REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    areas JSONB NOT NULL,
    deadline TIMESTAMPTZ NOT NULL,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_health_check_team_id
    ON thunderdome.health_check(team_id, created_at DESC);

-- user_id limits members to one response each, it is never returned with the answers
CREATE TABLE IF NOT EXISTS thunderdome.health_check_response (
    health_check_id UUID NOT NULL REFERENCES thunderdome.health_check(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES thunderdome.users(id) ON DELETE CASCADE,
    answers JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (health_check_id, user_id)
);

INSERT INTO thunderdome.health_check_template (name, description, is_public, areas) VALUES
(
    'Squad Health Check',
    'Based on the Spotify Squad Health Check model',
    true,
    '[
        {"name": "Delivering value", "healthy": "We deliver great things that we are proud of and our stakeholders are happy.", "unhealthy": "What we deliver is poor and our stakeholders are unhappy."},
        {"name": "Easy to release", "healthy": "Releasing is simple, safe, painless and mostly automated.", "unhealthy": "Releasing is risky, painful, slow and full of manual work."},
        {"name": "Health of codebase", "healthy": "Our code is clean, easy to read and well tested.", "unhealthy": "Our code is hard to work with and technical debt is out of control."},
        {"name": "Fun", "healthy": "We love coming to work and enjoy working together.", "unhealthy": "Work is boring or draining."},
        {"name": "Learning", "healthy": "We are learning interesting things all the time.", "unhealthy": "We never have time to learn anything."},
        {"name": "Mission", "healthy": "We know why we are here and are excited about it.", "unhealthy": "We have no clear picture of why we are here or where we are going."},
        {"name": "Pawns or players", "healthy": "We control our own destiny and decide what to build and how.", "unhealthy": "We are told what to do and have no say."},
        {"name": "Speed", "healthy": "We get things done quickly without waiting or delays.", "unhealthy": "We never seem to finish anything and keep getting stuck or interrupted."},
        {"name": "Suitable process", "healthy": "Our way of working fits us.", "unhealthy": "Our way of working gets in our way."},
        {"name": "Support", "healthy": "We get help and support when we ask for it.", "unhealthy": "We keep getting stuck because we cannot get the help we ask for."}
    ]'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.health_check_response;
DROP TABLE IF EXISTS thunderdome.health_check;
DROP TABLE IF EXISTS thunderdome.health_check_template;

-- +goose StatementEnd
