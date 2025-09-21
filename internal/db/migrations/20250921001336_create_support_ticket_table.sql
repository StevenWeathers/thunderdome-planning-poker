-- +goose Up
-- +goose StatementBegin

-- Support form table schema
CREATE TABLE thunderdome.support_ticket (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    inquiry TEXT NOT NULL,
    assigned_to UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    notes TEXT,
    resolved_at TIMESTAMPTZ,
    resolved_by UUID REFERENCES thunderdome.users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for common queries
CREATE INDEX idx_support_ticket_user_id ON thunderdome.support_ticket(user_id);
CREATE INDEX idx_support_ticket_assigned_to ON thunderdome.support_ticket(assigned_to);
CREATE INDEX idx_support_ticket_resolved_at ON thunderdome.support_ticket(resolved_at);
CREATE INDEX idx_support_ticket_resolved_by ON thunderdome.support_ticket(resolved_by);
CREATE INDEX idx_support_ticket_created_at ON thunderdome.support_ticket(created_at);
CREATE INDEX idx_support_ticket_email ON thunderdome.support_ticket(email);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS thunderdome.support_ticket;
-- +goose StatementEnd
