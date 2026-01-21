# DEPLOYMENT_GUIDE.md

## Purpose

This document provides **step-by-step deployment procedures** for Bangladesh Jatiotabadi Jubodal Management System (BJDMS) to the production server.

**Critical Constraint**: ALL deployment happens on **remote server only** (grayhawks.com). Local machines are code-only environments.

---

## Deployment Architecture

**Target Server**:
- **Host**: grayhawks.com
- **OS**: Debian 13
- **Access**: SSH only (`ssh root@grayhawks.com -p 25920`)
- **Runtime**: Docker + Docker Compose
- **Project Path**: `/opt/bjdms`

---

## Pre-Deployment Checklist

### 1. Code Preparation (Local)
- [ ] All code changes committed to git repository
- [ ] All tests passing (run remotely before deployment)
- [ ] Documentation updated
- [ ] Environment variables documented (not committed)
- [ ] Database migration scripts prepared (if schema changes)

### 2. Server Preparation (Remote)
- [ ] SSH access verified
- [ ] Docker and Docker Compose installed
- [ ] Required directories exist: `/opt/bjdms`, `/opt/bjdms/volumes`
- [ ] SSL certificates valid
- [ ] Backup completed (see BACKUP_RECOVERY.md)

---

## Deployment Steps

### Step 1: Connect to Server

```bash
# From local machine
ssh root@grayhawks.com -p 25920
```

**Verification**: `pwd` should show `/root` or similar.

---

### Step 2: Navigate to Project Directory

```bash
cd /opt/bjdms
```

---

### Step 3: Pull Latest Code

```bash
# If using git
git pull origin main

# If using rsync from local (alternative)
# From local machine:
rsync -avz -e "ssh -p 25920" \
  --exclude='.git' --exclude='node_modules' --exclude='*.env' \
  /path/to/local/bjdms/ root@grayhawks.com:/opt/bjdms/
```

**Verification**: `git log -1` shows latest commit.

---

### Step 4: Update Environment Variables

```bash
nano /opt/bjdms/.env
```

**Required Variables** (example):
```env
# Database
DATABASE_URL=postgresql://bjdms_user:STRONG_PASSWORD@postgres:5432/bjdms_db
POSTGRES_USER=bjdms_user
POSTGRES_PASSWORD=STRONG_PASSWORD
POSTGRES_DB=bjdms_db

# Redis
REDIS_URL=redis://redis:6379

# JWT
JWT_SECRET=RANDOM_256_BIT_SECRET
JWT_REFRESH_SECRET=ANOTHER_RANDOM_SECRET

# App
NODE_ENV=production
NEXT_PUBLIC_API_URL=https://grayhawks.com/api/v1

# OpenSearch  
OPENSEARCH_URL=http://opensearch:9200

# S3/MinIO (file storage)
S3_ENDPOINT=http://minio:9000
S3_ACCESS_KEY=ACCESS_KEY
S3_SECRET_KEY=SECRET_KEY
S3_BUCKET=bjdms-files

# SMS Gateway (Bangladesh)
SMS_API_KEY=YOUR_SMS_API_KEY
SMS_SENDER_ID=BJDMS

# Email (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=APP_PASSWORD
```

**Security**: Rotate secrets quarterly, never commit `.env` to git.

---

### Step 5: Run Database Migrations

**If schema changes exist**:

```bash
# Using Go migrations (example)
docker-compose run --rm api go run cmd/migrate/main.go up

# OR using Python/Alembic (example)
docker-compose run --rm api alembic upgrade head
```

**Verification**: Check `schema_migrations` table or equivalent.

**Rollback Plan**: Always test migrations on staging first. If migration fails, rollback:
```bash
go run cmd/migrate/main.go down 1
# or
alembic downgrade -1
```

---

### Step 6: Build Docker Images

```bash
docker-compose build --no-cache
```

**Note**: Rebuilds all images fresh. Use `--no-cache` to ensure latest dependencies.

**Expected Output**: "Successfully built..." for each service.

---

### Step 7: Stop Current Services (if updating)

```bash
docker-compose down
```

**Warning**: This causes downtime. For zero-downtime, see Blue-Green Deployment section.

---

### Step 8: Start Services

```bash
docker-compose up -d
```

**Flags**:
- `-d`: Detached mode (runs in background)

**Expected Output**: All services show "done".

---

### Step 9: Verify Services

```bash
# Check all containers running
docker-compose ps

# Should show:
# bjdms_api_1        running
# bjdms_postgres_1   running
# bjdms_redis_1      running
# bjdms_opensearch_1 running
# bjdms_analytics_1  running
# bjdms_web_1        running
```

**Check Logs**:
```bash
# API logs
docker-compose logs -f api

# Analytics logs
docker-compose logs -f analytics
```

---

### Step 10: Health Check

```bash
# From server
curl http://localhost:3000/health

# Expected response:
# {"status": "healthy", "database": "connected", "redis": "connected"}
```

**Public Health Check**:
```bash
# From local machine or browser
curl https://grayhawks.com/api/v1/health
```

---

### Step 11: Smoke Tests

**Test Authentication**:
```bash
curl -X POST https://grayhawks.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone": "+880XXXXXXXXX", "password": "test_password"}'
```

**Expected**: JWT tokens returned.

**Test Database**:
```bash
# Log into postgres container
docker-compose exec postgres psql -U bjdms_user -d bjdms_db

# Run query
SELECT COUNT(*) FROM users;
```

**Test Redis**:
```bash
docker-compose exec redis redis-cli PING
# Expected: PONG
```

---

### Step 12: Monitor Logs

```bash
# Watch logs for errors (first 30 minutes)
docker-compose logs -f --tail=100
```

**Look for**:
- No error messages
- Successful database connections
- Successful Redis connections
- No authentication failures

---

## Zero-Downtime Deployment (Blue-Green)

**For production with active users**:

### 1. Run New Stack in Parallel

```bash
# Copy docker-compose.yml to docker-compose.green.yml
cp docker-compose.yml docker-compose.green.yml

# Edit ports to avoid conflicts (e.g., 3001 instead of 3000)
nano docker-compose.green.yml

# Start green environment
docker-compose -f docker-compose.green.yml up -d
```

### 2. Health Check Green

```bash
curl http://localhost:3001/health
```

### 3. Update NGINX to Route to Green

```bash
nano /etc/nginx/sites-available/bjdms

# Change upstream from port 3000 to 3001
upstream bjdms_api {
    server localhost:3001;
}

# Reload NGINX
nginx -t && systemctl reload nginx
```

### 4. Stop Blue Environment

```bash
docker-compose down
```

### 5. Cleanup

Rename green to blue for next deployment.

---

## Rollback Procedure

**If deployment fails**:

### 1. Revert Code

```bash
git reset --hard HEAD~1  # Go back one commit
# or
git checkout {previous_commit_hash}
```

### 2. Revert Database Migrations

```bash
# Rollback last migration
go run cmd/migrate/main.go down 1
```

### 3. Rebuild and Restart

```bash
docker-compose build
docker-compose up -d
```

### 4. Verify

```bash
curl https://grayhawks.com/api/v1/health
```

**Rollback Time**: Target < 5 minutes.

---

## Post-Deployment Tasks

### 1. Verify Functionality

- [ ] Login with test account
- [ ] Submit test activity
- [ ] Create test complaint
- [ ] Search functionality working
- [ ] File upload working

### 2. Monitor Metrics

- [ ] CPU usage normal (<70%)
- [ ] Memory usage normal (<80%)
- [ ] Disk space available (>20% free)
- [ ] Response times normal (<500ms)

### 3. Backup Verification

- [ ] Automated backup ran successfully
- [ ] Backup file size reasonable

### 4. Notify Users (if downtime occurred)

SMS/Dashboard notification:
"System maintenance complete. BJDMS is now fully operational. Thank you for your patience."

---

## Deployment Frequency

**Recommended Schedule**:
- **Hotfixes**: As needed (security, critical bugs)
- **Minor Updates**: Weekly (Friday evening, low traffic)
- **Major Updates**: Monthly (first Saturday, with user notification)

**Maintenance Window**: Friday 11 PM - Saturday 1 AM Bangladesh Time (low usage period)

---

## Environment-Specific Configurations

### Development (Local Simulation - Remote Only)

Even for "development", use remote server:

```bash
# On grayhawks.com, create dev environment
mkdir /opt/bjdms-dev
cd /opt/bjdms-dev

# Use dev-specific docker-compose
docker-compose -f docker-compose.dev.yml up -d
```

**Access**: `https://grayhawks.com/dev/`

---

### Staging

**Purpose**: Test before production deployment.

**Location**: Same server, different directory (`/opt/bjdms-staging`)

**Domain**: `https://grayhawks.com/staging/` or subdomain

**Process**:
1. Deploy to staging first
2. Run full test suite
3. Manual QA testing
4. If successful, deploy to production

---

### Production

**Location**: `/opt/bjdms`

**Domain**: `https://grayhawks.com` or `https://bjdms.grayhawks.com`

**Process**: Follow main deployment steps above.

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs {service_name}

# Common issues:
# - Port already in use: Change port in docker-compose.yml
# - Missing environment variable: Check .env file
# - Database connection failed: Verify DATABASE_URL
```

---

### Database Migration Failed

```bash
# Rollback migration
go run cmd/migrate/main.go down 1

# Fix migration file
nano migrations/XXX_migration_name.sql

# Re-run
go run cmd/migrate/main.go up
```

---

### Out of Disk Space

```bash
# Check disk usage
df -h

# Clean Docker (removes unused images/volumes)
docker system prune -a

# Clean old logs
find /opt/bjdms/logs -type f -mtime +30 -delete
```

---

## Security Considerations

- **SSH Key Only**: Disable password authentication
- **Firewall**: Only ports 443, 25920 open
- **SSL Certificates**: Auto-renew (Certbot cron job)
- **Secrets Rotation**: Quarterly JWT secret rotation
- **Dependency Updates**: Monthly security patches

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
