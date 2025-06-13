CREATE TABLE quiz_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Participation tracking
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE, -- When they started answering
    submitted_at TIMESTAMP WITH TIME ZONE, -- When they submitted final answers
    disconnected_at TIMESTAMP WITH TIME ZONE, -- If they disconnected early

    -- Performance metrics
    total_score DECIMAL(5,2) DEFAULT 0.00,
    max_possible_score DECIMAL(5,2) DEFAULT 0.00,
    completion_percentage DECIMAL(5,2) DEFAULT 0.00,
    total_time_seconds INTEGER DEFAULT 0,

    -- Answer analytics (aggregated)
    questions_answered INTEGER DEFAULT 0,
    questions_correct INTEGER DEFAULT 0,
    questions_skipped INTEGER DEFAULT 0,
    average_response_time_ms INTEGER DEFAULT 0,
    fastest_response_time_ms INTEGER,
    slowest_response_time_ms INTEGER,

    -- Status tracking
    status VARCHAR(20) DEFAULT 'joined' CHECK (
        status IN ('joined', 'active', 'completed', 'disconnected', 'abandoned')
    ),

    UNIQUE(session_id, user_id) -- One participation per user per session
);

CREATE INDEX idx_participants_session ON quiz_participants(session_id);
CREATE INDEX idx_participants_user ON quiz_participants(user_id);
CREATE INDEX idx_participants_status ON quiz_participants(status);
CREATE INDEX idx_participants_score ON quiz_participants(total_score DESC);
CREATE INDEX idx_participants_completion ON quiz_participants(completion_percentage DESC);
