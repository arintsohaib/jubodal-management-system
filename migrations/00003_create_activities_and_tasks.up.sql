-- 1. Activities
CREATE TYPE activity_category AS ENUM ('political', 'social', 'organizational', 'protest', 'other');

CREATE TABLE IF NOT EXISTS activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) NOT NULL,
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL,
    committee_id UUID REFERENCES committees(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    category activity_category DEFAULT 'organizational',
    activity_date TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_activity_jurisdiction ON activities(jurisdiction_id);
CREATE INDEX idx_activity_user ON activities(user_id);

-- 2. Activity Proofs (Files)
CREATE TABLE IF NOT EXISTS activity_proofs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    activity_id UUID REFERENCES activities(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL, -- S3/MinIO path
    file_type VARCHAR(50),
    file_size INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 3. Tasks
CREATE TYPE task_status AS ENUM ('pending', 'in_progress', 'completed', 'verified', 'cancelled');

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES users(id) NOT NULL,
    assignee_id UUID REFERENCES users(id), -- Specific user
    committee_id UUID REFERENCES committees(id), -- Or whole committee
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status DEFAULT 'pending',
    priority INTEGER DEFAULT 2, -- 1: Critical, 2: High, 3: Medium, 4: Low
    due_date TIMESTAMP,
    completed_at TIMESTAMP,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_task_assignee ON tasks(assignee_id);
CREATE INDEX idx_task_committee ON tasks(committee_id);
CREATE INDEX idx_task_jurisdiction ON tasks(jurisdiction_id);

-- 4. Events
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL,
    organizer_id UUID REFERENCES users(id) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    location TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP,
    is_public BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- 5. Event Attendance
CREATE TABLE IF NOT EXISTS event_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) NOT NULL,
    attended_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(event_id, user_id)
);
