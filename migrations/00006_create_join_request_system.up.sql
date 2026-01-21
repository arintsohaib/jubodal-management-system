-- 1. Join Request Status Type
CREATE TYPE join_request_status AS ENUM ('pending', 'under_review', 'approved', 'rejected', 'completed');

-- 2. Join Requests Table
CREATE TABLE IF NOT EXISTS join_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    full_name_bn VARCHAR(255),
    phone VARCHAR(20) UNIQUE NOT NULL,
    nid VARCHAR(20) UNIQUE NOT NULL,
    date_of_birth DATE NOT NULL,
    gender VARCHAR(10),
    blood_group VARCHAR(5),
    occupation VARCHAR(100),
    address TEXT NOT NULL,
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL, -- Target jurisdiction to join
    applied_at TIMESTAMP DEFAULT NOW(),
    status join_request_status DEFAULT 'pending',
    referred_by_id UUID REFERENCES users(id), -- Optional referral
    rejection_reason TEXT,
    processed_by_id UUID REFERENCES users(id), -- Who gave final approval
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 3. Join Request Audit Logs
CREATE TABLE IF NOT EXISTS join_request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID REFERENCES join_requests(id) ON DELETE CASCADE,
    actor_id UUID REFERENCES users(id),
    action VARCHAR(50) NOT NULL, -- 'status_change', 'comment'
    old_status join_request_status,
    new_status join_request_status,
    note TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_join_request_jurisdiction ON join_requests(jurisdiction_id);
CREATE INDEX idx_join_request_status ON join_requests(status);
CREATE INDEX idx_join_request_phone ON join_requests(phone);
CREATE INDEX idx_join_request_nid ON join_requests(nid);
