CREATE TYPE quiz_session_status AS ENUM (
    'created',
    'scheduled',
    'active',
    'paused',
    'ended',
    'completed',
    'archived',
    'cancelled'
);
CREATE TABLE quiz_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    classroom_id UUID NOT NULL REFERENCES classrooms(id) ON DELETE CASCADE,
    conducted_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,

    session_code VARCHAR(10) UNIQUE NOT NULL, -- Short code for students to join
    session_name VARCHAR(255),
    status quiz_session_status DEFAULT 'created',
    current_question_index INTEGER DEFAULT 0,
    scheduled_start_at TIMESTAMP WITH TIME ZONE,
    actual_start_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    total_duration_seconds INTEGER,

    total_participants INTEGER DEFAULT 0,
    active_participants INTEGER DEFAULT 0,
    completed_participants INTEGER DEFAULT 0,

    average_score DECIMAL(5,2) DEFAULT 0.00,
    highest_score DECIMAL(5,2) DEFAULT 0.00,
    lowest_score DECIMAL(5,2) DEFAULT 0.00,

    total_interactions INTEGER DEFAULT 0, -- number of all user actions
    peak_concurrent_users INTEGER DEFAULT 0,

    settings JSONB DEFAULT '{}', -- Session-specific settings

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_sessions_quiz ON quiz_sessions(quiz_id);
CREATE INDEX idx_sessions_classroom ON quiz_sessions(classroom_id);
CREATE INDEX idx_sessions_code ON quiz_sessions(session_code);
CREATE INDEX idx_sessions_status ON quiz_sessions(status);
CREATE INDEX idx_sessions_start_time ON quiz_sessions(actual_start_at DESC);
