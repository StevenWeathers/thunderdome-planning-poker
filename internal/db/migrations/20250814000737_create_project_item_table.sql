-- +goose Up
-- +goose StatementBegin

-- Create color enum type
CREATE TYPE thunderdome.color_type AS ENUM (
    'red',
    'orange', 
    'amber',
    'yellow',
    'lime',
    'green',
    'emerald',
    'teal',
    'cyan',
    'sky',
    'blue',
    'indigo',
    'violet',
    'purple',
    'fuchsia',
    'pink',
    'rose'
);

-- Item Types table for global, org, department, and team scoped types
CREATE TABLE thunderdome.item_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    type_key VARCHAR(10) NOT NULL,
    description TEXT,
    color thunderdome.color_type, -- color from enum type
    is_active BOOLEAN DEFAULT true,
    organization_id UUID REFERENCES thunderdome.organization(id), -- NULL for global types
    department_id UUID REFERENCES thunderdome.organization_department(id), -- NULL for global/org types
    team_id UUID REFERENCES thunderdome.team(id), -- NULL for global/org/dept types
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(type_key, organization_id, department_id, team_id), -- ensures unique type_key within scope
    -- Ensure proper hierarchy: team belongs to dept, dept belongs to org
    CHECK (
        (organization_id IS NULL AND department_id IS NULL AND team_id IS NULL) OR -- Global
        (organization_id IS NOT NULL AND department_id IS NULL AND team_id IS NULL) OR -- Org-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NULL) OR -- Dept-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NOT NULL) -- Team-level
    ),
    -- Ensure type_key contains only lowercase letters, underscores, and dashes
    CHECK (type_key ~ '^[a-z_-]+$')
);

-- Item Statuses table for global, org, department, and team scoped statuses
CREATE TABLE thunderdome.item_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    status_key VARCHAR(20) NOT NULL,
    description TEXT,
    color thunderdome.color_type, -- color from enum type
    is_initial BOOLEAN DEFAULT false, -- marks default initial status
    is_final BOOLEAN DEFAULT false, -- marks completion/closure statuses
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    organization_id UUID REFERENCES thunderdome.organization(id), -- NULL for global statuses
    department_id UUID REFERENCES thunderdome.organization_department(id), -- NULL for global/org statuses
    team_id UUID REFERENCES thunderdome.team(id), -- NULL for global/org/dept statuses
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(status_key, organization_id, department_id, team_id), -- ensures unique status_key within scope
    -- Ensure proper hierarchy: team belongs to dept, dept belongs to org
    CHECK (
        (organization_id IS NULL AND department_id IS NULL AND team_id IS NULL) OR -- Global
        (organization_id IS NOT NULL AND department_id IS NULL AND team_id IS NULL) OR -- Org-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NULL) OR -- Dept-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NOT NULL) -- Team-level
    )
);

-- Item Priorities table for global, org, department, and team scoped priorities
CREATE TABLE thunderdome.item_priority (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL,
    priority_key VARCHAR(20) NOT NULL,
    description TEXT,
    color thunderdome.color_type, -- color from enum type
    priority_level INTEGER NOT NULL, -- numeric value for sorting (lower = higher priority)
    is_active BOOLEAN DEFAULT true,
    organization_id UUID REFERENCES thunderdome.organization(id), -- NULL for global priorities
    department_id UUID REFERENCES thunderdome.organization_department(id), -- NULL for global/org priorities
    team_id UUID REFERENCES thunderdome.team(id), -- NULL for global/org/dept priorities
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(priority_key, organization_id, department_id, team_id), -- ensures unique priority_key within scope
    UNIQUE(priority_level, organization_id, department_id, team_id), -- ensures unique priority levels within scope
    -- Ensure proper hierarchy: team belongs to dept, dept belongs to org
    CHECK (
        (organization_id IS NULL AND department_id IS NULL AND team_id IS NULL) OR -- Global
        (organization_id IS NOT NULL AND department_id IS NULL AND team_id IS NULL) OR -- Org-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NULL) OR -- Dept-level
        (organization_id IS NOT NULL AND department_id IS NOT NULL AND team_id IS NOT NULL) -- Team-level
    )
);

-- Project Item table
CREATE TABLE thunderdome.project_item (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES thunderdome.project(id),
    parent_id UUID REFERENCES thunderdome.project_item(id),
    item_key VARCHAR(20) NOT NULL, -- project_key + incremental number (e.g. PROJ-123)
    title VARCHAR(255) NOT NULL,
    description TEXT,
    type_id UUID NOT NULL REFERENCES thunderdome.item_type(id),
    status_id UUID NOT NULL REFERENCES thunderdome.item_status(id),
    priority_id UUID REFERENCES thunderdome.item_priority(id), -- nullable for items without priority
    story_points VARCHAR(8),
    rank VARCHAR(16) NOT NULL,
    start_date TIMESTAMPTZ,
    end_date TIMESTAMPTZ,
    created_by UUID NOT NULL REFERENCES thunderdome.users(id),
    external_reference_id VARCHAR(100),
    external_reference_link TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, item_key) -- ensures unique item keys within project scope
);

-- Insert global default item types
INSERT INTO thunderdome.item_type (name, type_key, description, organization_id, department_id, team_id) VALUES
('epic', 'epic', 'Large body of work that can be broken down into stories', NULL, NULL, NULL),
('story', 'story', 'User story or feature request', NULL, NULL, NULL),
('spike', 'spike', 'Research or investigation task', NULL, NULL, NULL),
('bug', 'bug', 'Software defect or issue', NULL, NULL, NULL),
('task', 'task', 'General work item', NULL, NULL, NULL),
('sub_task', 'sub-task', 'Subtask of another item', NULL, NULL, NULL);

-- Insert global default item statuses
INSERT INTO thunderdome.item_status (name, status_key, description, is_initial, is_final, sort_order, organization_id, department_id, team_id) VALUES
('not_started', 'not-started', 'Work has not begun', true, false, 1, NULL, NULL, NULL),
('backlog', 'backlog', 'In the backlog awaiting prioritization', false, false, 2, NULL, NULL, NULL),
('to_do', 'to-do', 'Ready to be worked on', false, false, 3, NULL, NULL, NULL),
('in_progress', 'in-progress', 'Currently being worked on', false, false, 4, NULL, NULL, NULL),
('in_review', 'in-review', 'Under review or testing', false, false, 5, NULL, NULL, NULL),
('blocked', 'blocked', 'Blocked by external dependency', false, false, 6, NULL, NULL, NULL),
('done', 'done', 'Work completed successfully', false, true, 7, NULL, NULL, NULL),
('completed', 'completed', 'Work completed and verified', false, true, 8, NULL, NULL, NULL),
('canceled', 'canceled', 'Work was canceled', false, true, 9, NULL, NULL, NULL),
('wont_do', 'wont-do', 'Decided not to do this work', false, true, 10, NULL, NULL, NULL),
('duplicate', 'duplicate', 'Duplicate of another item', false, true, 11, NULL, NULL, NULL);

-- Insert global default item priorities
INSERT INTO thunderdome.item_priority (name, priority_key, description, priority_level, organization_id, department_id, team_id) VALUES
('lowest', 'lowest', 'Lowest priority', 7, NULL, NULL, NULL),
('low', 'low', 'Low priority', 6, NULL, NULL, NULL),
('medium', 'medium', 'Medium priority', 5, NULL, NULL, NULL),
('high', 'high', 'High priority', 4, NULL, NULL, NULL),
('highest', 'highest', 'Highest priority', 3, NULL, NULL, NULL),
('urgent', 'urgent', 'Urgent - needs immediate attention', 2, NULL, NULL, NULL),
('blocker', 'blocker', 'Critical blocker', 1, NULL, NULL, NULL);

-- Indexes for performance
CREATE INDEX idx_item_type_org_dept_team ON thunderdome.item_type(organization_id, department_id, team_id);
CREATE INDEX idx_item_status_org_dept_team ON thunderdome.item_status(organization_id, department_id, team_id);
CREATE INDEX idx_item_priority_org_dept_team ON thunderdome.item_priority(organization_id, department_id, team_id);
CREATE INDEX idx_item_status_sort_order ON thunderdome.item_status(sort_order);
CREATE INDEX idx_item_priority_level ON thunderdome.item_priority(priority_level);
CREATE INDEX idx_project_item_key ON thunderdome.project_item(item_key);
CREATE INDEX idx_project_item_project_id ON thunderdome.project_item(project_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

- Drop indexes
DROP INDEX IF EXISTS thunderdome.idx_item_priority_level;
DROP INDEX IF EXISTS thunderdome.idx_item_status_sort_order;
DROP INDEX IF EXISTS thunderdome.idx_item_priority_org_dept_team;
DROP INDEX IF EXISTS thunderdome.idx_item_status_org_dept_team;
DROP INDEX IF EXISTS thunderdome.idx_item_type_org_dept_team;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS thunderdome.project_item;
DROP TABLE IF EXISTS thunderdome.item_priority;
DROP TABLE IF EXISTS thunderdome.item_status;
DROP TABLE IF EXISTS thunderdome.item_type;

-- Drop enum type
DROP TYPE IF EXISTS thunderdome.color_type;

-- +goose StatementEnd
