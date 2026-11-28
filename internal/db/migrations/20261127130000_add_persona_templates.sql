-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS thunderdome.persona_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    personas JSONB NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT false,
    created_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    organization_id UUID REFERENCES thunderdome.organization(id) ON DELETE CASCADE,
    team_id UUID REFERENCES thunderdome.team(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT persona_template_check_scope CHECK (
        (is_public = true AND num_nonnulls(organization_id, team_id) = 0) OR
        (is_public = false AND num_nonnulls(organization_id, team_id) = 1)
    )
);

CREATE INDEX IF NOT EXISTS idx_persona_template_organization_id
    ON thunderdome.persona_template(organization_id);

CREATE INDEX IF NOT EXISTS idx_persona_template_team_id
    ON thunderdome.persona_template(team_id);

CREATE INDEX IF NOT EXISTS idx_persona_template_is_public
    ON thunderdome.persona_template(is_public);

INSERT INTO thunderdome.persona_template (name, description, is_public, personas) VALUES
(
    'Software as a Service',
    'Common users of a subscription software product',
    true,
    '[
        {"name": "New User", "role": "Evaluator", "description": "Recently signed up and wants to reach a first success quickly."},
        {"name": "Power User", "role": "Daily user", "description": "Relies on the product every day and cares about speed and shortcuts."},
        {"name": "Account Administrator", "role": "Admin", "description": "Manages members, permissions and billing for their organization."}
    ]'
),
(
    'Online Store',
    'Common users of an e-commerce site',
    true,
    '[
        {"name": "First-time Shopper", "role": "Visitor", "description": "Browsing to see whether the store has what they need and can be trusted."},
        {"name": "Returning Customer", "role": "Customer", "description": "Wants to find, reorder and track purchases with little effort."},
        {"name": "Store Manager", "role": "Staff", "description": "Keeps the catalog, stock and orders up to date."}
    ]'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS thunderdome.persona_template;

-- +goose StatementEnd
