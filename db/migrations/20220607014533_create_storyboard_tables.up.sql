CREATE TABLE IF NOT EXISTS storyboard (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(256),
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    join_code VARCHAR(128),
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    color_legend JSONB DEFAULT '[{"color":"gray","legend":""},{"color":"red","legend":""},{"color":"orange","legend":""},{"color":"yellow","legend":""},{"color":"green","legend":""},{"color":"teal","legend":""},{"color":"blue","legend":""},{"color":"indigo","legend":""},{"color":"purple","legend":""},{"color":"pink","legend":""}]'::JSONB
);

CREATE TABLE IF NOT EXISTS storyboard_goal (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    name VARCHAR(256),
    sort_order INTEGER,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(storyboard_id, sort_order)
);

CREATE TABLE IF NOT EXISTS storyboard_column (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    goal_id UUID REFERENCES storyboard_goal(id) ON DELETE CASCADE,
    name VARCHAR(256),
    sort_order INTEGER,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(goal_id, sort_order)
);

CREATE TABLE IF NOT EXISTS storyboard_story (
    id UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    goal_id UUID REFERENCES storyboard_goal(id) ON DELETE CASCADE,
    column_id UUID REFERENCES storyboard_column(id) ON DELETE CASCADE,
    name VARCHAR(256),
    color VARCHAR(32) DEFAULT 'gray',
    content TEXT,
    sort_order INTEGER,
    points INTEGER,
    closed BOOL DEFAULT false,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(column_id, sort_order)
);

CREATE TABLE IF NOT EXISTS storyboard_user (
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    active BOOL DEFAULT false,
    abandoned BOOL DEFAULT false,
    PRIMARY KEY (storyboard_id, user_id)
);

CREATE TABLE IF NOT EXISTS storyboard_persona (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,
    role VARCHAR(256),
    description TEXT,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(storyboard_id, name)
);

CREATE TABLE IF NOT EXISTS storyboard_story_comment (
    id UUID NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    story_id UUID REFERENCES storyboard_story(id) ON DELETE CASCADE,
    comment TEXT,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS team_storyboard (
    team_id UUID REFERENCES team(id) ON DELETE CASCADE,
    storyboard_id UUID REFERENCES storyboard(id) ON DELETE CASCADE,
    created_date TIMESTAMPTZ DEFAULT NOW(),
    updated_date TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (team_id, storyboard_id)
);