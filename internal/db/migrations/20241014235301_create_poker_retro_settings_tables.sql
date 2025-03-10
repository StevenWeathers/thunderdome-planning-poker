-- +goose Up
-- +goose StatementBegin

CREATE TABLE thunderdome.poker_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    department_id UUID REFERENCES thunderdome.organization_department(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    auto_finish_voting BOOLEAN DEFAULT true,
    point_average_rounding VARCHAR(5) DEFAULT 'ceil',
    hide_voter_identity BOOLEAN DEFAULT false,
    estimation_scale_id UUID REFERENCES thunderdome.estimation_scale(id),
    join_code TEXT,
    facilitator_code TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_id_not_null CHECK (
        num_nonnulls(organization_id, department_id, team_id) = 1
    ),
    CONSTRAINT ps_unique_org_setting UNIQUE (organization_id),
    CONSTRAINT ps_unique_dept_setting UNIQUE (department_id),
    CONSTRAINT ps_unique_team_setting UNIQUE (team_id)
);

CREATE INDEX ON thunderdome.poker_settings (organization_id);
CREATE INDEX ON thunderdome.poker_settings (department_id);
CREATE INDEX ON thunderdome.poker_settings (team_id);

CREATE TABLE thunderdome.retro_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    department_id UUID REFERENCES thunderdome.organization_department(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    max_votes INT2 DEFAULT 3,
    allow_multiple_votes BOOLEAN DEFAULT false,
    brainstorm_visibility VARCHAR(12) DEFAULT 'visible',
    phase_time_limit_min INT2 DEFAULT 0,
    phase_auto_advance BOOLEAN DEFAULT false,
    allow_cumulative_voting BOOLEAN DEFAULT false,
    template_id UUID REFERENCES thunderdome.retro_template(id),
    join_code TEXT,
    facilitator_code TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT check_id_not_null CHECK (
        num_nonnulls(organization_id, department_id, team_id) = 1
    ),
    CONSTRAINT rs_unique_org_setting UNIQUE (organization_id),
    CONSTRAINT rs_unique_dept_setting UNIQUE (department_id),
    CONSTRAINT rs_unique_team_setting UNIQUE (team_id)
);

CREATE INDEX ON thunderdome.retro_settings (organization_id);
CREATE INDEX ON thunderdome.retro_settings (department_id);
CREATE INDEX ON thunderdome.retro_settings (team_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.poker_settings;
DROP TABLE IF EXISTS thunderdome.retro_settings;

-- +goose StatementEnd