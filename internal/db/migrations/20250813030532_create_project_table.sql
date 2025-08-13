-- +goose Up
-- +goose StatementBegin

-- Main project table
CREATE TABLE thunderdome.project (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_key VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    department_id UUID REFERENCES thunderdome.organization_department(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure at least one association exists
    CONSTRAINT check_project_association CHECK (
        organization_id IS NOT NULL OR
        department_id IS NOT NULL OR
        team_id IS NOT NULL
    ),
    
    -- Ensure key format is valid (alphanumeric, uppercase)
    CONSTRAINT check_project_key_format CHECK (
        project_key ~ '^[A-Z0-9]{2,10}$'
    )
);

-- Unique indexes ensuring key uniqueness within each scope
CREATE UNIQUE INDEX project_org_key_unique
ON thunderdome.project (organization_id, project_key)
WHERE organization_id IS NOT NULL;

CREATE UNIQUE INDEX project_dept_key_unique
ON thunderdome.project (department_id, project_key)
WHERE department_id IS NOT NULL;

CREATE UNIQUE INDEX project_team_key_unique
ON thunderdome.project (team_id, project_key)
WHERE team_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.project;

-- +goose StatementEnd