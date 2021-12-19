CREATE TABLE IF NOT EXISTS retro (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(256),
    format VARCHAR(32) NOT NULL DEFAULT 'worked_improve_question', -- start_stop_continue
    phase VARCHAR(16) NOT NULL DEFAULT 'intro',
    join_code VARCHAR(128),
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS retro_action (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    completed BOOL DEFAULT false,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS retro_group (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    name VARCHAR(128),
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS retro_group_vote (
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    group_id UUID REFERENCES retro_group(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (retro_id, user_id, group_id)
);

CREATE TABLE IF NOT EXISTS retro_item (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    group_id UUID REFERENCES retro_group(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    type VARCHAR(16) NOT NULL, -- worked, improve, question, start, stop, continue
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS retro_user (
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    active BOOL DEFAULT false,
    abandoned BOOL DEFAULT false,
    spectator BOOL DEFAULT false,
    PRIMARY KEY (retro_id, user_id)
);

CREATE TABLE IF NOT EXISTS team_retro (
    team_id UUID REFERENCES team(id) ON DELETE CASCADE,
    retro_id UUID REFERENCES retro(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (team_id, retro_id)
);
