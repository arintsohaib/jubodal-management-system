# RELEASE_MANAGEMENT.md

## Purpose

Defines **release planning, versioning, and deployment procedures** for BJDMS.

**Goal**: Controlled, predictable releases with minimal user disruption.

---

## Versioning Strategy

**Semantic Versioning** (SemVer): `MAJOR.MINOR.PATCH`

**Examples**:
- `1.0.0`: Initial stable release
- `1.1.0`: New feature (backward compatible)
- `1.1.1`: Bug fix
- `2.0.0`: Breaking change (API change, schema migration requiring downtime)

---

### Version Components

- **MAJOR** (1.x.x): Breaking changes
  - API endpoint removed/renamed
  - Database schema change requiring migration WITH downtime
  - Permission model change

- **MINOR** (x.1.x): New features (backward compatible)
  - New API endpoint
  - New table added (no downtime)
  - New UI feature

- **PATCH** (x.x.1): Bug fixes
  - Security patches
  - Bug fixes
  - Performance improvements (no new features)

---

## Release Cycle

**Scheduled Releases**:
- **Major**: Annually (e.g., v2.0.0 on January 1)
- **Minor**: Monthly (first Saturday of month)
- **Patch**: As needed (security/critical bugs)

**Release Candidate**: 1 week before release
- Tag: `v1.1.0-rc1`
- Deploy to staging
- QA testing
- If issues found → `v1.1.0-rc2`

---

## Release Process

### 1. Planning Phase (2 weeks before)

**Release Manager**:
- Create GitHub milestone: `v1.1.0`
- Identify features/fixes for release
- Review/prioritize issues
- Assign to developers

---

### 2. Code Freeze (1 week before)

**No new features** after code freeze. Only:
- Bug fixes
- Documentation
- Tests

**Branch**: `release/v1.1.0`
```bash
git checkout -b release/v1.1.0 main
git push origin release/v1.1.0
```

---

### 3. Release Candidate

**Tag**:
```bash
git tag -a v1.1.0-rc1 -m "Release candidate 1 for v1.1.0"
git push origin v1.1.0-rc1
```

**Deploy to Staging**:
```bash
ssh root@arint.win -p 25920
cd /opt/bjdms-staging
git checkout v1.1.0-rc1
./scripts/deploy.sh
```

**QA Testing**:
- Manual testing of new features
- Regression testing
- Performance testing
- Security review

**If bugs found**:
- Fix in `release/v1.1.0` branch
- Tag `v1.1.0-rc2`
- Re-test

---

### 4. Final Release

**When**: All tests pass, no critical bugs

**Tag**:
```bash
git checkout release/v1.1.0
git tag -a v1.1.0 -m "Version 1.1.0 - Municipality Support"
git push origin v1.1.0
```

**Merge to Main**:
```bash
git checkout main
git merge release/v1.1.0
git push origin main
```

---

### 5. Production Deployment

**Deployment Window**: Friday 11 PM - Saturday 1 AM (low traffic)

**Pre-Deployment**:
- [ ] Backup database
- [ ] Notify users (24 hours advance)
- [ ] Prepare rollback plan
- [ ] Ops team on standby

**Deploy**:
```bash
ssh root@arint.win -p 25920
cd /opt/bjdms
git checkout v1.1.0
./scripts/deploy.sh
```

**Post-Deployment**:
- [ ] Smoke tests
- [ ] Monitor logs for 1 hour
- [ ] Monitor metrics (CPU, memory, errors)
- [ ] Notify users (deployment complete)

---

### 6. Release Notes

**Publish**: GitHub Releases + Website announcement

**Template**:
```markdown
# v1.1.0 - Municipality Support

Released: 2026-02-01

## New Features
- Added Municipality level to jurisdiction hierarchy
- Anonymous complaint submission (no login required)
- Bengali language search support

## Improvements
- Faster committee member search (20% improvement)
- Better error messages for validation failures

## Bug Fixes
- Fixed phone number validation (#123)
- Fixed committee dissolution audit log (#145)

## Database Changes
- Added `is_urban` column to `jurisdictions` table
- Added `join_requests` table
- Migration required (zero downtime)

## Upgrade Instructions
1. Backup database
2. Pull latest code: `git pull origin v1.1.0`
3. Run migrations: `./scripts/deploy.sh`
4. Verify health check

## Breaking Changes
None

## Security Fixes
- Fixed SQL injection in complaint search (CVE-2026-0001)
```

---

## Hotfix Process

**For Critical Bugs in Production**

### 1. Create Hotfix Branch

```bash
git checkout -b hotfix/security-patch main
```

---

### 2. Fix & Test

- Implement fix
- Write test to prevent regression
- Test on staging

---

### 3. Release

**Tag**:
```bash
git tag -a v1.1.2 -m "Hotfix: Security patch for SQL injection"
git push origin v1.1.2
```

**Deploy Immediately** (no waiting for release window):
```bash
ssh root@arint.win -p 25920
cd /opt/bjdms
git checkout v1.1.2
./scripts/deploy.sh
```

---

### 4. Backport to Main

```bash
git checkout main
git merge hotfix/security-patch
git push origin main
```

---

## User Communication

### Pre-Release Announcement (24 hours before)

**SMS/Dashboard Notification**:
```
BJDMS will undergo maintenance on Saturday, Feb 1, 12:00 AM - 2:00 AM. 
New features: Municipality support, anonymous complaints. 
Minimal downtime expected. Thank you for your patience.
```

---

### Release Announcement

**Dashboard Banner**:
```
🎉 BJDMS v1.1.0 Released!
New: Municipality level support, anonymous complaints, Bengali search.
Learn more: [Release Notes]
```

---

### Post-Release Follow-Up

**48 hours after**:
- Monitor support requests
- Track bug reports
- Gather user feedback

---

## Rollback Plan

**If Critical Issue Found**:

### 1. Identify Problem

- Error spike in logs
- User reports
- Monitoring alerts

---

### 2. Decide: Fix Forward or Rollback?

**Fix Forward** (preferred if quick fix available):
- Deploy hotfix within 30 minutes
- Less disruptive

**Rollback** (if issue unclear or fix complex):
- Restore previous version
- Investigate offline

---

### 3. Execute Rollback

```bash
ssh root@arint.win -p 25920
cd /opt/bjdms
git checkout v1.0.1  # Previous stable version
./scripts/deploy.sh
```

**If Database Migration**:
- Rollback migration:
  ```bash
  go run cmd/migrate/main.go down 1
  ```
- Restore from backup if needed

---

### 4. Communicate

**Immediate Notification**:
```
BJDMS has been temporarily rolled back to v1.0.1 due to a technical issue. 
We are investigating and will provide updates shortly. We apologize for the inconvenience.
```

---

## Release Checklist

### Pre-Release

- [ ] All features complete
- [ ] All tests passing
- [ ] Documentation updated
- [ ] Security review completed
- [ ] Performance testing completed
- [ ] Database migration tested
- [ ] Rollback plan documented
- [ ] Release notes drafted
- [ ] User notification sent (24h before)

### During Release

- [ ] Database backup completed
- [ ] Code deployed
- [ ] Migrations run
- [ ] Health check passed
- [ ] Smoke tests passed
- [ ] Monitoring active

### Post-Release

- [ ] Release notes published
- [ ] User notification sent (release complete)
- [ ] Monitoring for 1 hour
- [ ] Team debriefing scheduled (next day)

---

এই ডকুমেন্ট ভাঙলে সিস্টেম ভাঙবে।
এই ডকুমেন্ট ঠিক থাকলে – যত বড়ই হোক – সিস্টেম স্থিতিশীল থাকবে।
