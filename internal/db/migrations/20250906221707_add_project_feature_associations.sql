
-- +goose Up
-- +goose StatementBegin
CREATE TABLE thunderdome.project_retro (
    retro_id UUID NOT NULL REFERENCES thunderdome.retro(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES thunderdome.project(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (retro_id, project_id)
);

CREATE TABLE thunderdome.project_poker (
    poker_id UUID NOT NULL REFERENCES thunderdome.poker(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES thunderdome.project(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (poker_id, project_id)
);

CREATE TABLE thunderdome.project_storyboard (
    storyboard_id UUID NOT NULL REFERENCES thunderdome.storyboard(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES thunderdome.project(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (storyboard_id, project_id)
);

CREATE INDEX idx_project_retro_retro_id ON thunderdome.project_retro(retro_id);
CREATE INDEX idx_project_retro_project_id ON thunderdome.project_retro(project_id);
CREATE INDEX idx_project_poker_poker_id ON thunderdome.project_poker(poker_id);
CREATE INDEX idx_project_poker_project_id ON thunderdome.project_poker(project_id);
CREATE INDEX idx_project_storyboard_storyboard_id ON thunderdome.project_storyboard(storyboard_id);
CREATE INDEX idx_project_storyboard_project_id ON thunderdome.project_storyboard(project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE thunderdome.project_retro;
DROP TABLE thunderdome.project_poker;
DROP TABLE thunderdome.project_storyboard;
-- +goose StatementEnd