CREATE TABLE classrooms (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL REFERENCES schools(id),
    name VARCHAR(255) NOT NULL,
    grade_level VARCHAR(50),
    subject VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
