-- 1. Complaint Status Type
CREATE TYPE complaint_status AS ENUM ('received', 'under_review', 'action_taken', 'closed', 'rejected');

-- 2. Complaints Table
CREATE TABLE IF NOT EXISTS complaints (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tracking_id VARCHAR(20) UNIQUE NOT NULL, -- Human readable tracking (e.g. C-2026-X7Y)
    user_id UUID REFERENCES users(id), -- Nullable for anonymous
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL, -- Target jurisdiction
    is_anonymous BOOLEAN DEFAULT TRUE,
    complainant_name VARCHAR(255), -- Null if anonymous
    complainant_contact VARCHAR(255), -- Null if anonymous
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    status complaint_status DEFAULT 'received',
    anonymous_ip_hash VARCHAR(64), -- SHA256 of IP for spam control
    assigned_to_id UUID REFERENCES users(id), -- Handling official
    resolution_notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_complaint_jurisdiction ON complaints(jurisdiction_id);
CREATE INDEX idx_complaint_tracking ON complaints(tracking_id);
CREATE INDEX idx_complaint_status ON complaints(status);

-- 3. Complaint Evidence
CREATE TABLE IF NOT EXISTS complaint_evidence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    complaint_id UUID REFERENCES complaints(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    file_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

-- 4. Complaint Comments/Audit Trail
CREATE TABLE IF NOT EXISTS complaint_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    complaint_id UUID REFERENCES complaints(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id), -- Null for system
    action VARCHAR(50) NOT NULL, -- 'status_change', 'comment', 'assigned'
    old_status complaint_status,
    new_status complaint_status,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
