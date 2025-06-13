CREATE TYPE user_role AS ENUM ('student', 'teacher');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255),
    name VARCHAR(255),
    role user_role NOT NULL,
    school_id UUID REFERENCES schools(id) ON DELETE CASCADE,

    -- Analytics tracking fields
    first_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    total_quiz_sessions INTEGER DEFAULT 0,
    total_questions_answered INTEGER DEFAULT 0,
    average_response_time_ms INTEGER DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_school_role ON users(school_id, role);
CREATE INDEX idx_users_last_seen ON users(last_seen_at DESC);
