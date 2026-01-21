-- 1. Transaction Types and Categories
CREATE TYPE transaction_type AS ENUM ('income', 'expense', 'transfer');

CREATE TABLE IF NOT EXISTS finance_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    name_bn VARCHAR(255) NOT NULL,
    type transaction_type NOT NULL,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Initial Categories
INSERT INTO finance_categories (name, name_bn, type, is_system) VALUES
('Donation', 'দান', 'income', true),
('Membership Fee', 'সদস্যপদ ফি', 'income', true),
('Event Sponsor', 'অনুষ্ঠান স্পনসর', 'income', false),
('Office Rent', 'অফিস ভাড়া', 'expense', true),
('Utility Bill', 'ইউটিলিটি বিল', 'expense', true),
('Printing & Stationary', 'মুদ্রণ এবং স্টেশনারি', 'expense', false),
('Transport', 'পরিবহন', 'expense', false),
('Grant/Aid', 'অনুদান/সহায়তা', 'expense', false);

-- 2. Financial Transactions
CREATE TABLE IF NOT EXISTS finance_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jurisdiction_id UUID REFERENCES jurisdictions(id) NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL, -- The person recording it
    category_id UUID REFERENCES finance_categories(id) NOT NULL,
    type transaction_type NOT NULL,
    amount DECIMAL(15, 2) NOT NULL CHECK (amount > 0),
    description TEXT,
    reference_no VARCHAR(100), -- Physical receipt no
    transaction_date DATE DEFAULT CURRENT_DATE,
    evidence_path TEXT, -- Link to receipt image
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Immutability check: No Updates or Deletes on transactions
CREATE OR REPLACE FUNCTION protect_transactions()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'Financial transactions are immutable and cannot be modified or deleted.';
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_protect_finance_transactions
BEFORE UPDATE OR DELETE ON finance_transactions
FOR EACH STATEMENT EXECUTE FUNCTION protect_transactions();

-- 3. Jurisdiction Balances (Cached for performance)
CREATE TABLE IF NOT EXISTS finance_balances (
    jurisdiction_id UUID PRIMARY KEY REFERENCES jurisdictions(id),
    total_income DECIMAL(15, 2) DEFAULT 0,
    total_expense DECIMAL(15, 2) DEFAULT 0,
    current_balance DECIMAL(15, 2) DEFAULT 0,
    last_updated_at TIMESTAMP DEFAULT NOW()
);

-- 4. Trigger to update balances automatically
CREATE OR REPLACE FUNCTION update_jurisdiction_balance()
RETURNS TRIGGER AS $$
BEGIN
    -- Initialize balance record if not exists
    INSERT INTO finance_balances (jurisdiction_id)
    VALUES (NEW.jurisdiction_id)
    ON CONFLICT (jurisdiction_id) DO NOTHING;

    IF NEW.type = 'income' THEN
        UPDATE finance_balances 
        SET total_income = total_income + NEW.amount,
            current_balance = current_balance + NEW.amount,
            last_updated_at = NOW()
        WHERE jurisdiction_id = NEW.jurisdiction_id;
    ELSIF NEW.type = 'expense' THEN
        UPDATE finance_balances 
        SET total_expense = total_expense + NEW.amount,
            current_balance = current_balance - NEW.amount,
            last_updated_at = NOW()
        WHERE jurisdiction_id = NEW.jurisdiction_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_finance_balance
AFTER INSERT ON finance_transactions
FOR EACH ROW EXECUTE FUNCTION update_jurisdiction_balance();

-- Indexes
CREATE INDEX idx_finance_jurisdiction ON finance_transactions(jurisdiction_id);
CREATE INDEX idx_finance_category ON finance_transactions(category_id);
CREATE INDEX idx_finance_date ON finance_transactions(transaction_date);
