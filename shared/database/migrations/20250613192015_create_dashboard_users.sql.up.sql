CREATE TABLE dashboard_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) DEFAULT 'analyst' CHECK (role IN ('viewer', 'analyst', 'admin')),

    -- Access control
    school_access UUID[], -- school IDs they can access (NULL = all schools)
    permissions JSONB DEFAULT '{}', -- Specific permissions

    is_active BOOLEAN DEFAULT TRUE,
    last_login_at TIMESTAMP WITH TIME ZONE,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_dashboard_users_username ON dashboard_users(username);
CREATE INDEX idx_dashboard_users_active ON dashboard_users(is_active);
