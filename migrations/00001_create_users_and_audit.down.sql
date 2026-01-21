-- Drop audit log trigger and function
DROP TRIGGER IF EXISTS audit_immutable_trigger ON audit_logs;
DROP FUNCTION IF EXISTS prevent_audit_modification();

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS users;
