-- 1. Jurisdiction Levels
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE IF NOT EXISTS jurisdiction_levels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- Central, Division, District, Upazila, Municipality, Union, Ward
    rank INTEGER UNIQUE NOT NULL,     -- 1 to 7
    description TEXT
);

INSERT INTO jurisdiction_levels (name, rank, description) VALUES
('Central', 1, 'National HQ'),
('Division', 2, 'Regional (8 Divisions)'),
('District', 3, 'District Level (64 Districts)'),
('Upazila', 4, 'Rural Administrative Unit'),
('Municipality', 5, 'Urban Administrative Unit (Pourashava)'),
('Union', 6, 'Rural local unit'),
('Ward', 7, 'Lowest administrative unit');

-- 2. Jurisdictions
CREATE TABLE IF NOT EXISTS jurisdictions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    level_id INTEGER REFERENCES jurisdiction_levels(id) NOT NULL,
    parent_id UUID REFERENCES jurisdictions(id), -- NULL for Central
    name VARCHAR(255) NOT NULL,
    name_bn VARCHAR(255),
    is_urban BOOLEAN DEFAULT FALSE,
    population INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_jurisdiction_parent ON jurisdictions(parent_id);
CREATE INDEX idx_jurisdiction_level ON jurisdictions(level_id);

-- 3. Positions (Targeted 16 Positions)
CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    name_bn VARCHAR(100),
    rank INTEGER NOT NULL, -- Low = High Priority (e.g., 1 = President)
    committee_type VARCHAR(20) DEFAULT 'Both', -- Full, Convener, Both
    description TEXT
);

INSERT INTO positions (name, name_bn, rank, committee_type) VALUES
('President', 'সভাপতি', 1, 'Full'),
('General Secretary', 'সাধারণ সম্পাদক', 2, 'Full'),
('Convener', 'আহ্বায়ক', 1, 'Convener'),
('Member Secretary', 'সদস্য সচিব', 2, 'Convener'),
('Senior Vice President', 'সিনিয়র সহ-সভাপতি', 3, 'Full'),
('Vice President', 'সহ-সভাপতি', 4, 'Full'),
('Joint General Secretary', 'যুগ্ম সাধারণ সম্পাদক', 5, 'Full'),
('Assistant General Secretary', 'সহ-সাধারণ সম্পাদক', 6, 'Full'),
('Organizational Secretary', 'সাংগঠনিক সম্পাদক', 7, 'Full'),
('Treasurer', 'কোষাধ্যক্ষ', 8, 'Both'),
('Office Secretary', 'দপ্তর সম্পাদক', 9, 'Full'),
('Publicity Secretary', 'প্রচার সম্পাদক', 10, 'Full'),
('Member', 'সদস্য', 100, 'Both');

-- 4. Committees
CREATE TYPE committee_status AS ENUM ('proposed', 'active', 'dissolved', 'expired');
CREATE TYPE committee_type AS ENUM ('full', 'convener');

CREATE TABLE IF NOT EXISTS committees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL,
    type committee_type NOT NULL,
    status committee_status DEFAULT 'proposed',
    formed_at TIMESTAMP,
    expires_at TIMESTAMP,
    approved_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    CONSTRAINT one_active_committee_per_jurisdiction EXCLUDE USING gist (jurisdiction_id WITH =, status WITH =) WHERE (status = 'active')
);

-- Note: EXCLUDE constraint requires btree_gist extension in Postgres

-- 5. Committee Members
CREATE TABLE IF NOT EXISTS committee_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    committee_id UUID REFERENCES committees(id) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    position_id INTEGER REFERENCES positions(id) NOT NULL,
    joined_at TIMESTAMP DEFAULT NOW(),
    ended_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(committee_id, user_id, position_id)
);

CREATE INDEX idx_member_user ON committee_members(user_id);
CREATE INDEX idx_member_committee ON committee_members(committee_id);

-- 6. User Verification (Verification link to committee)
-- Adding a column to users to link them to their current active position for quick lookups
ALTER TABLE users ADD COLUMN IF NOT EXISTS current_committee_id UUID REFERENCES committees(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS current_position_id INTEGER REFERENCES positions(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS jurisdiction_id UUID REFERENCES jurisdictions(id);
