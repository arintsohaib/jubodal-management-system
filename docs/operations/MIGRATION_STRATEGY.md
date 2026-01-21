# MIGRATION_STRATEGY.md

## Purpose

Defines **database schema migration strategy** for BJDMS to evolve the system safely without data loss or downtime.

**Principle**: Never break existing functionality. All migrations must be backward-compatible or have clear rollback paths.

---

## Migration Tools

**Recommended**:
- **Go**: `golang-migrate/migrate` or `goose`
- **Python**: `Alembic` (if using Python/FastAPI services)
- **Node.js**: Prisma Migrate, TypeORM, or Knex.js

**Choice depends on primary backend language**. Document chosen tool in this file once decided.

---

## Migration File Structure

**Directory**: `/opt/bjdms/migrations/`

**Naming Convention**: `{TIMESTAMP}_{description}.{up|down}.sql`

**Example**:
```
migrations/
├── 20260115120000_create_users_table.up.sql
├── 20260115120000_create_users_table.down.sql
├── 20260118143000_add_municipality_level.up.sql
├── 20260118143000_add_municipality_level.down.sql
└── 20260121100000_add_join_requests_table.up.sql
└── 20260121100000_add_join_requests_table.down.sql
```

---

## Migration Workflow

### 1. Create Migration

**Manual (SQL)**:
```bash
# Generate timestamp
TIMESTAMP=$(date +%Y%m%d%H%M%S)

# Create migration files
touch migrations/${TIMESTAMP}_add_full_name_bn_to_users.up.sql
touch migrations/${TIMESTAMP}_add_full_name_bn_to_users.down.sql
```

**UP Migration** (`add_full_name_bn_to_users.up.sql`):
```sql
ALTER TABLE users ADD COLUMN full_name_bn VARCHAR(255);
CREATE INDEX idx_users_full_name_bn ON users(full_name_bn);
```

**DOWN Migration** (`.down.sql`):
```sql
DROP INDEX IF EXISTS idx_users_full_name_bn;
ALTER TABLE users DROP COLUMN IF EXISTS full_name_bn;
```

---

### 2. Test Locally (on Remote Server)

**Never test migrations on production first!**

```bash
# On grayhawks.com, in test environment
cd /opt/bjdms-test

# Apply migration
migrate -path migrations -database "postgresql://user:pass@localhost/bjdms_test_db" up

# Verify
psql -U user -d bjdms_test_db -c "\d users"

# Test rollback
migrate -path migrations -database "postgresql://user:pass@localhost/bjdms_test_db" down 1

# Verify rollback worked
psql -U user -d bjdms_test_db -c "\d users"
```

---

### 3. Code Review

Migration reviewed for:
- [ ] Correct SQL syntax
- [ ] Backward compatibility (won't break running app)
- [ ] Rollback works
- [ ] Indexes added for new columns (if queried)
- [ ] No data loss
- [ ] Tested in test environment

---

### 4. Deploy to Production

**Pre-Deployment**:
- [ ] Backup database (see BACKUP_RECOVERY.md)
- [ ] Notify users if downtime expected (rare)

**Execute**:
```bash
ssh root@grayhawks.com -p 25920
cd /opt/bjdms

# Apply migration
migrate -path migrations -database "$DATABASE_URL" up

# OR using Docker
docker-compose run --rm api migrate -path migrations -database "$DATABASE_URL" up
```

**Verification**:
```bash
# Check schema
docker-compose exec postgres psql -U bjdms_user -d bjdms_db -c "\d+ users"

# Check application still works
curl https://grayhawks.com/api/v1/health
```

---

### 5. Monitor Post-Migration

- Check application logs for errors
- Monitor query performance (new columns/indexes)
- Test affected features manually

---

## Migration Types

### 1. Additive Migrations (Safe, No Downtime)

**Examples**:
- Add new column (nullable or with default)
- Add new table
- Add index
- Add check constraint (non-blocking)

**Deployment**: Can be applied while application is running.

**SQL Example**:
```sql
-- Safe: New column with default
ALTER TABLE users ADD COLUMN email_verified BOOLEAN DEFAULT FALSE;

-- Safe: New table
CREATE TABLE notifications (id UUID PRIMARY KEY, ...);

-- Safe: Add index (CONCURRENTLY to avoid table lock)
CREATE INDEX CONCURRENTLY idx_users_phone ON users(phone);
```

---

### 2. Destructive Migrations (Requires Caution)

**Examples**:
- Drop column
- Drop table
-Rename column (breaks old code)
- Change column type

**Deployment**: Requires multi-step approach (see below).

---

### 3. Data Migrations (Modifies Data)

**Examples**:
- Populate new column from existing data
- Update enum values
- Normalize data

**Safe Approach**:
```sql
-- Step 1: Add column (additive, safe)
ALTER TABLE users ADD COLUMN jurisdiction_type VARCHAR(10);

-- Step 2: Populate data (in batches to avoid long lock)
UPDATE users
SET jurisdiction_type = CASE
  WHEN jurisdiction_id IN (SELECT id FROM jurisdictions WHERE is_urban = true) THEN 'urban'
  ELSE 'rural'
END
WHERE jurisdiction_type IS NULL
LIMIT 1000;  -- Batch size, repeat until all updated

-- Step 3: Make NOT NULL (only after all data populated)
ALTER TABLE users ALTER COLUMN jurisdiction_type SET NOT NULL;
```

---

## Safe Migration Patterns

### Pattern 1: Column Rename (3-Step)

**Goal**: Rename `full_name` to `name`

**Step 1**: Add new column (deploy code supports both)
```sql
ALTER TABLE users ADD COLUMN name VARCHAR(255);
UPDATE users SET name = full_name WHERE name IS NULL;
```

**Code**: Update to read/write both `full_name` and `name`.

**Step 2**: Backfill old data
```sql
UPDATE users SET name = full_name WHERE name IS NULL;
```

**Step 3**: Drop old column (after code only uses `name`)
```sql
ALTER TABLE users DROP COLUMN full_name;
```

**Timeline**: 1 week between each step.

---

### Pattern 2: Column Type Change

**Goal**: Change `phone VARCHAR(20)` to `phone VARCHAR(30)`

**Safe** (expansion):
```sql
ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(30);
```

**Unsafe** (contraction, e.g., VARCHAR(30) → VARCHAR(15)):
1. Verify no data exceeds new limit:
   ```sql
   SELECT COUNT(*) FROM users WHERE LENGTH(phone) > 15;
   ```
2. If count = 0, safe to proceed:
   ```sql
   ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(15);
   ```

---

### Pattern 3: Add NOT NULL Constraint

**Unsafe** (if nulls exist):
```sql
ALTER TABLE users ALTER COLUMN email SET NOT NULL;  -- Fails if nulls exist
```

**Safe**:
```sql
-- Step 1: Add with default
ALTER TABLE users ADD COLUMN email VARCHAR(255) DEFAULT 'noemail@example.com';

-- Step 2: Populate real values
UPDATE users SET email = ... WHERE email = 'noemail@example.com';

-- Step 3: Add NOT NULL
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
```

---

## Rollback Strategy

### When to Rollback

- Migration causes application errors
- Migration takes too long (>5 minutes)
- Data corruption detected

### How to Rollback

```bash
# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Restart application (may have cached old schema)
docker-compose restart api
```

### Rollback Testing

ALWAYS test rollback before deploying migration:
```bash
# In test environment
migrate up    # Apply migration
# Test app...
migrate down 1  # Rollback
# Test app again (should still work)
```

---

## Long-Running Migrations

**Problem**: Large table alterations can lock table for minutes/hours.

**Solutions**:

### 1. Create Index CONCURRENTLY (PostgreSQL)
```sql
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
```
No table lock, safe for production.

### 2. Use pt-online-schema-change (MySQL) or pg_repack (PostgreSQL)

Creates shadow table, copies data incrementally, swaps table.

### 3. Background Jobs
For data migrations, use background job:
```sql
-- Create column
ALTER TABLE users ADD COLUMN calculated_field INTEGER;

-- App code: Worker job populates in batches
UPDATE users SET calculated_field = ... WHERE id IN (SELECT id FROM users WHERE calculated_field IS NULL LIMIT 1000);
```

---

## Migration Checklist

Before deploying migration:
- [ ] UP and DOWN migrations written
- [ ] Tested in test environment (up + down)
- [ ] Backward compatible or multi-step approach planned
- [ ] Database backup created
- [ ] Code changes deployed (if needed)
- [ ] Rollback plan documented
- [ ] Monitoring in place for post-migration issues

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
