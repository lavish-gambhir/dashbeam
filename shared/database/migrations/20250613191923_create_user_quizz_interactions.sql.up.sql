CREATE TABLE user_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id UUID REFERENCES quiz_sessions(id) ON DELETE CASCADE, -- NULL for non-quiz interactions

    -- Interaction details
    interaction_type VARCHAR(50) NOT NULL, -- 'question_view', 'answer_select', 'app_focus', etc.
    interaction_target VARCHAR(100), -- What they interacted with
    interaction_data JSONB, -- Additional interaction details

    -- Context
    app_type VARCHAR(10) NOT NULL CHECK (app_type IN ('white', 'note')), -- Which app
    screen_name VARCHAR(100), -- Current screen/page
    classroom_id UUID REFERENCES classrooms(id) ON DELETE SET NULL,

    -- Timing
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    session_duration_ms INTEGER, -- How long they've been in current session

    -- Technical metadata
    device_info JSONB, -- Device type, OS version, etc.
    app_version VARCHAR(20),
    network_type VARCHAR(20), -- wifi, cellular, etc.

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_interactions_user ON user_interactions(user_id);
CREATE INDEX idx_interactions_session ON user_interactions(session_id);
CREATE INDEX idx_interactions_type ON user_interactions(interaction_type);
CREATE INDEX idx_interactions_timestamp ON user_interactions(timestamp DESC);
CREATE INDEX idx_interactions_app_type ON user_interactions(app_type);
CREATE INDEX idx_interactions_user_timestamp ON user_interactions(user_id, timestamp DESC);
