CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    classroom_id UUID REFERENCES classrooms(id) ON DELETE CASCADE,

    -- Quiz configuration
    total_questions INTEGER NOT NULL,
    time_limit_seconds INTEGER, -- NULL means no time limit
    max_attempts INTEGER DEFAULT 1,
    shuffle_questions BOOLEAN DEFAULT FALSE,
    shuffle_options BOOLEAN DEFAULT FALSE,

    -- Analytics fields
    total_sessions INTEGER DEFAULT 0,
    total_participants INTEGER DEFAULT 0,
    average_score DECIMAL(5,2) DEFAULT 0.00,
    average_completion_time_seconds INTEGER DEFAULT 0,

    -- Metadata
    tags TEXT[], -- For categorization and filtering
    difficulty_level VARCHAR(20), -- easy, medium, hard
    subject VARCHAR(100),

    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('draft', 'active', 'archived')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_quizzes_classroom ON quizzes(classroom_id);
CREATE INDEX idx_quizzes_creator ON quizzes(created_by_user_id);
CREATE INDEX idx_quizzes_status ON quizzes(status);
CREATE INDEX idx_quizzes_subject ON quizzes(subject);
CREATE INDEX idx_quizzes_difficulty ON quizzes(difficulty_level);
CREATE INDEX idx_quizzes_tags ON quizzes USING gin(tags);
