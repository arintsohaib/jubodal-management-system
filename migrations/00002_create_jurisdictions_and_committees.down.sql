DROP TABLE IF EXISTS committee_members;
DROP TABLE IF EXISTS committees;
DROP TYPE IF EXISTS committee_type;
DROP TYPE IF EXISTS committee_status;
DROP TABLE IF EXISTS positions;
DROP TABLE IF EXISTS jurisdictions;
DROP TABLE IF EXISTS jurisdiction_levels;

ALTER TABLE users DROP COLUMN IF EXISTS current_committee_id;
ALTER TABLE users DROP COLUMN IF EXISTS current_position_id;
ALTER TABLE users DROP COLUMN IF EXISTS jurisdiction_id;
