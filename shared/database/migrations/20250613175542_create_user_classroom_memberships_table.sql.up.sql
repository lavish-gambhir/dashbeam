CREATE TABLE user_classroom_memberships (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    classroom_id UUID NOT NULL REFERENCES classrooms(id) ON DELETE CASCADE,
    role user_role NOT NULL DEFAULT 'student',
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'pending')),

    UNIQUE(user_id, classroom_id) -- One membership per user per classroom
);

CREATE INDEX idx_memberships_user ON user_classroom_memberships(user_id);
CREATE INDEX idx_memberships_classroom ON user_classroom_memberships(classroom_id);
CREATE INDEX idx_memberships_classroom_role ON user_classroom_memberships(classroom_id, role);
CREATE INDEX idx_memberships_status ON user_classroom_memberships(status);
